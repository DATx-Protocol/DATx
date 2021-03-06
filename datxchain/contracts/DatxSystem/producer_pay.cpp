#include "DatxSystem.hpp"
#include <ctime>
#include <cmath>
#include <DatxToken/DatxToken.hpp>

namespace datxossystem {

   const int64_t  min_pervote_daily_pay = 100'0000;
   const int64_t  min_activated_stake   = 1'125'000'000'0000;
   double   continuous_rate       = 0.04879;          // 5% annual rate
   const double   perblock_rate         = 0.0025;           // 0.25%
   const double   standby_rate          = 0.0075;           // 0.75%
   const uint32_t blocks_per_year       = 52*7*24*2*3600;   // half seconds per year
   const uint32_t seconds_per_year      = 52*7*24*3600;
   const uint32_t blocks_per_day        = 2 * 24 * 3600;
   const uint32_t blocks_per_hour       = 2 * 3600;
   const uint64_t useconds_per_day      = 24 * 3600 * uint64_t(1000000);
   const uint64_t useconds_per_year     = seconds_per_year*1000000ll;
   //发行货币时间
   const time_t token_start_time = 1483200000;  //senconds,2017-1-1 0:0:0

   void system_contract::onblock( block_timestamp timestamp, account_name producer ) {
      using namespace datxos;

      require_auth(N(datxos));

      /** until activated stake crosses this threshold no new rewards are paid */
      if( _gstate.total_activated_stake < min_activated_stake )
         return;

      if( _gstate.last_pervote_bucket_fill == 0 )  /// start the presses
         _gstate.last_pervote_bucket_fill = current_time();


      /**
       * At startup the initial producer may not be one that is registered / elected
       * and therefore there may be no producer object for them.
       */
      auto prod = _producers.find(producer);
      if ( prod != _producers.end() ) {
         _gstate.total_unpaid_blocks++;
         _producers.modify( prod, 0, [&](auto& p ) {
               p.unpaid_blocks++;
         });
      }

      /// only update block producers once every minute, block_timestamp is in half seconds
      if( timestamp.slot - _gstate.last_producer_schedule_update.slot > 120 ) {
         update_elected_producers( timestamp );

         if( (timestamp.slot - _gstate.last_name_close.slot) > blocks_per_day ) {
            name_bid_table bids(_self,_self);
            auto idx = bids.get_index<N(highbid)>();
            auto highest = idx.begin();
            if( highest != idx.end() &&
                highest->high_bid > 0 &&
                highest->last_bid_time < (current_time() - useconds_per_day) &&
                _gstate.thresh_activated_stake_time > 0 &&
                (current_time() - _gstate.thresh_activated_stake_time) > 14 * useconds_per_day ) {
                   _gstate.last_name_close = timestamp;
                   idx.modify( highest, 0, [&]( auto& b ){
                         b.high_bid = -b.high_bid;
               });
            }
         }
      }
   }

   using namespace datxos;
   void system_contract::claimrewards( const account_name& owner ) {
      require_auth(owner);
      //auto producer = _producers.find(producer);
      
      const auto& prod = _producers.get( owner );
    
       
      datxos_assert( prod.active(), "producer does not have an active key" );
      

      datxos_assert( _gstate.total_activated_stake >= min_activated_stake,
                    "cannot claim rewards until the chain is activated (at least 15% of all tokens participate in voting)" );

      auto ct = current_time();

      datxos_assert( ct - prod.last_claim_time > useconds_per_day, "already claimed rewards within past day" );

      const asset token_supply   = token( N(datxos.token)).get_supply(symbol_type(system_token_symbol).name() );
      const auto usecs_since_last_fill = ct - _gstate.last_pervote_bucket_fill;

      if( usecs_since_last_fill > 0 && _gstate.last_pervote_bucket_fill > 0 ) {
        time_t time_now = time(0);
        time_t time_diff = time_now -token_start_time;
	double time_diff_year = time_diff / (3600.000 * 24 * 356);
	uint32_t diff_int= floor(time_diff_year);
         //连续增发的数量 
        continuous_rate=log(1.0+1/(70.0+diff_int));
         auto new_tokens = static_cast<int64_t>( (continuous_rate * double(token_supply.amount) * double(usecs_since_last_fill)) / double(useconds_per_year) );

         auto to_producers       = new_tokens / 5;  //20%
         auto to_savings         = new_tokens - to_producers;
         auto to_per_block_pay   = to_producers *0.35;  //35% block
         auto to_per_vote_pay    = to_producers *0.40;  //40% vote
         auto to_per_verify_pay    = to_producers *0.25;  //25% verify

         INLINE_ACTION_SENDER(datxos::token, issue)( N(datxos.token), {{N(datxos),N(active)}},
                                                    {N(datxos), asset(new_tokens), std::string("issue tokens for producer pay and savings")} );

         INLINE_ACTION_SENDER(datxos::token, transfer)( N(datxos.token), {N(datxos),N(active)},
                                                       { N(datxos), N(datxos.save), asset(to_savings), "unallocated inflation" } );

         INLINE_ACTION_SENDER(datxos::token, transfer)( N(datxos.token), {N(datxos),N(active)},
                                                       { N(datxos), N(datxos.bpay), asset(to_per_block_pay), "fund per-block bucket" } );

         INLINE_ACTION_SENDER(datxos::token, transfer)( N(datxos.token), {N(datxos),N(active)},
                                                       { N(datxos), N(datxos.vpay), asset(to_per_vote_pay), "fund per-vote bucket" } );
         auto verifier = _verifiers.find(owner);
         if ( verifier != _verifiers.end() ) {
	 INLINE_ACTION_SENDER(datxos::token, transfer)( N(datxos.token), {N(datxos),N(active)},
                                                       { N(datxos), N(datxos.veri), asset(to_per_verify_pay), "fund per-verify bucket" } );

	 INLINE_ACTION_SENDER(datxos::token, transfer)( N(datxos.token), {N(datxos.veri),N(active)},
                                                       { N(datxos.veri), owner, asset(to_per_verify_pay), std::string("producer verify pay") } );
         }
         _gstate.pervote_bucket  += to_per_vote_pay;
         _gstate.perblock_bucket += to_per_block_pay;

         _gstate.last_pervote_bucket_fill = ct;
      }

      int64_t producer_per_block_pay = 0;
      if( _gstate.total_unpaid_blocks > 0 ) {
         producer_per_block_pay = (_gstate.perblock_bucket * prod.unpaid_blocks) / _gstate.total_unpaid_blocks;
      }
      int64_t producer_per_vote_pay = 0;
      if( _gstate.total_producer_vote_weight > 0 ) {
         producer_per_vote_pay  = int64_t((_gstate.pervote_bucket * prod.total_votes ) / _gstate.total_producer_vote_weight);
      }
      if( producer_per_vote_pay < min_pervote_daily_pay ) {
         producer_per_vote_pay = 0;
      }
      _gstate.pervote_bucket      -= producer_per_vote_pay;
      _gstate.perblock_bucket     -= producer_per_block_pay;
      _gstate.total_unpaid_blocks -= prod.unpaid_blocks;

      _producers.modify( prod, 0, [&](auto& p) {
          p.last_claim_time = ct;
          p.unpaid_blocks = 0;
      });

      if( producer_per_block_pay > 0 ) {
         INLINE_ACTION_SENDER(datxos::token, transfer)( N(datxos.token), {N(datxos.bpay),N(active)},
                                                       { N(datxos.bpay), owner, asset(producer_per_block_pay), std::string("producer block pay") } );
      }
      if( producer_per_vote_pay > 0 ) {
         INLINE_ACTION_SENDER(datxos::token, transfer)( N(datxos.token), {N(datxos.vpay),N(active)},
                                                       { N(datxos.vpay), owner, asset(producer_per_vote_pay), std::string("producer vote pay") } );
      }  

     
      
      
      
   }

} //namespace datxossystem
