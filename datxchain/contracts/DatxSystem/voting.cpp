/**
 *  @file
 *  @copyright defined in DATX/LICENSE.txt
 *
**/


#include "DatxSystem.hpp"

#include <DatxosLib/DatxosLib.hpp>
#include <DatxosLib/crypto.h> 
#include <DatxosLib/print.hpp>
#include <DatxosLib/datastream.hpp>
#include <DatxosLib/serialize.hpp>
#include <DatxosLib/multi_index.hpp>
#include <DatxosLib/privileged.hpp>
#include <DatxosLib/singleton.hpp>
#include <DatxosLib/transaction.hpp>
#include <DatxToken/DatxToken.hpp>

#include <algorithm>
#include <cmath>
#include <ctime>
#include <string>


namespace datxossystem {
   using datxos::indexed_by;
   using datxos::const_mem_fun;
   using datxos::bytes;
   using datxos::print;
   using datxos::singleton;
   using datxos::transaction;

   /**
    *  This method will create a producer_config and producer_info object for 'producer'
    *
    *  @pre producer is not already registered
    *  @pre producer to register is an account
    *  @pre authority of producer to register
    *
    */
   void system_contract::regproducer( const account_name producer, const datxos::public_key& producer_key, const std::string& url,const std::string& url_verifier, uint16_t location ) {
      datxos_assert( url.size() < 512, "url too long" );
      datxos_assert( producer_key != datxos::public_key(), "public key should not be the default value" );
      require_auth( producer );

      auto prod = _producers.find( producer );

      if ( prod != _producers.end() ) {
         _producers.modify( prod, producer, [&]( producer_info& info ){
               info.producer_key = producer_key;
               info.is_active    = true;
               info.url          = url;
               info.location     = location;
            });
      } else {
         _producers.emplace( producer, [&]( producer_info& info ){
               info.owner         = producer;
               info.total_votes   = 0;
               info.producer_key  = producer_key;
               info.is_active     = true;
               info.url           = url;
               info.location      = location;
         });
      }
        //verifiers
      auto veri = _verifiers.find( producer );

      if ( veri != _verifiers.end() ) {
         _verifiers.modify( veri, producer, [&]( verifier_info& info ){
               info.verifier_key = producer_key;
               info.is_active    = true;
               info.url          = url_verifier;
               info.location     = location;
            });
      } else {
         _verifiers.emplace( producer, [&]( verifier_info& info ){
               info.owner         = producer;
               info.total_votes   = 0;
               info.verifier_key  = producer_key;
               info.is_active     = true;
               info.url           = url_verifier;
               info.location      = location;
         });
      } 
   }
   

   void system_contract::unregprod( const account_name producer ) {
      require_auth( producer );

      const auto& prod = _producers.get( producer, "producer not found" );

      _producers.modify( prod, 0, [&]( producer_info& info ){
            info.deactivate();
      });
   }

   void system_contract::update_elected_producers( block_timestamp block_time ) {
      _gstate.last_producer_schedule_update = block_time;

      auto idx = _producers.get_index<N(prototalvote)>();

      std::vector< std::pair<datxos::producer_key,uint16_t> > top_producers;
      top_producers.reserve(_maxbp);

      for ( auto it = idx.cbegin(); it != idx.cend() && top_producers.size() < _maxbp && 0 < it->total_votes && it->active(); ++it ) {
         top_producers.emplace_back( std::pair<datxos::producer_key,uint16_t>({{it->owner, it->producer_key}, it->location}) );
      }

      if ( top_producers.size() < _gstate.last_producer_schedule_size ) {
         return;
      }

      /// sort by producer name
      std::sort( top_producers.begin(), top_producers.end() );

      std::vector<datxos::producer_key> producers;

      producers.reserve(top_producers.size());
      for( const auto& item : top_producers )
         producers.push_back(item.first);

      bytes packed_schedule = pack(producers);

      if( set_proposed_producers( packed_schedule.data(),  packed_schedule.size() ) >= 0 ) {
         _gstate.last_producer_schedule_size = static_cast<decltype(_gstate.last_producer_schedule_size)>( top_producers.size() );
      }
      using namespace datxos;
      int tmp_time=0;
      
      //_gstate.last_verifier_schedule_update = block_time;
      if ((block_time.slot-_gstate.last_verifier_schedule_update.slot)%246857==0)
      {
          int count=0;
          tmp_time=1;
          for ( auto it = idx.cbegin(); it != idx.cend() && count < _maxbp && 0 < it->total_votes && it->active(); ++it ){

          auto prod = _verifiers.find( it->owner );
        if ( prod != _verifiers.end() ) {
           _verifiers.modify( prod, it->owner, [&]( verifier_info& info ){
                 info.verifier_key = it->producer_key;
                 info.is_active    = true;
                 info.url          = it->url;
                 info.location     = it->location;
              });
          } else {
           _verifiers.emplace( it->owner, [&]( verifier_info& info ){
                 info.owner         = it->owner;
                 info.total_votes   = 0;
                 info.verifier_key  = it->producer_key;
                 info.is_active     = true;
                 info.url           = it->url;
                 info.location      = it->location;
           });
          }
          count++;
          }//end for
      }
//	  time_t verify_start_time = 1483200000;
//      //require_auth(N(datxos));
//      time_t time_now = time(0);
//      time_t time_diff = time_now - verify_start_time;
//      if ( (time_diff) / (3600*24*356 / 2) == 0)
//      {
//        int count=0;

//        for ( auto it = idx.cbegin(); it != idx.cend() && count < _maxbp && 0 < it->total_votes && it->active(); ++it ){
	
//        auto prod = _verifiers.find( it->owner );
//	if ( prod != _verifiers.end() ) {
//         _verifiers.modify( prod, it->owner, [&]( verifier_info& info ){
//               info.verifier_key = it->producer_key;
//               info.is_active    = true;
//               info.url          = it->url;
//               info.location     = it->location;
//            });
//      	} else {
//         _verifiers.emplace( it->owner, [&]( verifier_info& info ){
//               info.owner         = it->owner;
//               info.total_votes   = 0;
//               info.verifier_key  = it->producer_key;
//               info.is_active     = true;
//               info.url           = it->url;
//               info.location      = it->location;
//         });
//      	}
//        count++;
//        }//end for
	
//      }//end if time()
     
   }

   double stake2vote( int64_t staked ) {
      /// TODO subtract 2080 brings the large numbers closer to this decade
      double weight = int64_t( (now() - (block_timestamp::block_timestamp_epoch / 1000)) / (seconds_per_day * 7) )  / double( 52 );
      return double(staked) * std::pow( 2, weight );
   }
   /**
    *  @pre producers must be sorted from lowest to highest and must be registered and active
    *  @pre if proxy is set then no producers can be voted for
    *  @pre if proxy is set then proxy account must exist and be registered as a proxy
    *  @pre every listed producer or proxy must have been previously registered
    *  @pre voter must authorize this action
    *  @pre voter must have previously staked some DATX for voting
    *  @pre voter->staked must be up to date
    *
    *  @post every producer previously voted for will have vote reduced by previous vote weight
    *  @post every producer newly voted for will have vote increased by new vote amount
    *  @post prior proxy will proxied_vote_weight decremented by previous vote weight
    *  @post new proxy will proxied_vote_weight incremented by new vote weight
    *
    *  If voting for a proxy, the producer votes will not change until the proxy updates their own vote.
    */
   void system_contract::voteproducer( const account_name voter_name, const account_name proxy, const std::vector<account_name>& producers ) {
      require_auth( voter_name );
      update_votes( voter_name, proxy, producers, true );
   }

   void system_contract::update_votes( const account_name voter_name, const account_name proxy, const std::vector<account_name>& producers, bool voting ) {
      //validate input
      if ( proxy ) {
         datxos_assert( producers.size() == 0, "cannot vote for producers and proxy at same time" );
         datxos_assert( voter_name != proxy, "cannot proxy to self" );
         require_recipient( proxy );
      } else {
         datxos_assert( producers.size() <= 30, "attempt to vote for too many producers" );
         for( size_t i = 1; i < producers.size(); ++i ) {
            datxos_assert( producers[i-1] < producers[i], "producer votes must be unique and sorted" );
         }
      }

      auto voter = _voters.find(voter_name);
      datxos_assert( voter != _voters.end(), "user must stake before they can vote" ); /// staking creates voter object
      datxos_assert( !proxy || !voter->is_proxy, "account registered as a proxy is not allowed to use a proxy" );

      /**
       * The first time someone votes we calculate and set last_vote_weight, since they cannot unstake until
       * after total_activated_stake hits threshold, we can use last_vote_weight to determine that this is
       * their first vote and should consider their stake activated.
       */
      if( voter->last_vote_weight <= 0.0 ) {
         _gstate.total_activated_stake += voter->staked;
         if( _gstate.total_activated_stake >= min_activated_stake && _gstate.thresh_activated_stake_time == 0 ) {
            _gstate.thresh_activated_stake_time = current_time();
         }
      }

      auto new_vote_weight = stake2vote( voter->staked );
      if( voter->is_proxy ) {
         new_vote_weight += voter->proxied_vote_weight;
      }

      boost::container::flat_map<account_name, pair<double, bool /*new*/> > producer_deltas;
      if ( voter->last_vote_weight > 0 ) {
         if( voter->proxy ) {
            auto old_proxy = _voters.find( voter->proxy );
            datxos_assert( old_proxy != _voters.end(), "old proxy not found" ); //data corruption
            _voters.modify( old_proxy, 0, [&]( auto& vp ) {
                  vp.proxied_vote_weight -= voter->last_vote_weight;
               });
            propagate_weight_change( *old_proxy );
         } else {
            for( const auto& p : voter->producers ) {
               auto& d = producer_deltas[p];
               d.first -= voter->last_vote_weight;
               d.second = false;
            }
         }
      }

      if( proxy ) {
         auto new_proxy = _voters.find( proxy );
         datxos_assert( new_proxy != _voters.end(), "invalid proxy specified" ); //if ( !voting ) { data corruption } else { wrong vote }
         datxos_assert( !voting || new_proxy->is_proxy, "proxy not found" );
         if ( new_vote_weight >= 0 ) {
            _voters.modify( new_proxy, 0, [&]( auto& vp ) {
                  vp.proxied_vote_weight += new_vote_weight;
               });
            propagate_weight_change( *new_proxy );
         }
      } else {
         if( new_vote_weight >= 0 ) {
            for( const auto& p : producers ) {
               auto& d = producer_deltas[p];
               d.first += new_vote_weight;
               d.second = true;
            }
         }
      }

      for( const auto& pd : producer_deltas ) {
         auto pitr = _producers.find( pd.first );
         if( pitr != _producers.end() ) {
            datxos_assert( !voting || pitr->active() || !pd.second.second /* not from new set */, "producer is not currently registered" );
            _producers.modify( pitr, 0, [&]( auto& p ) {
               p.total_votes += pd.second.first;
               if ( p.total_votes < 0 ) { // floating point arithmetics can give small negative numbers
                  p.total_votes = 0;
               }
               _gstate.total_producer_vote_weight += pd.second.first;
               //datxos_assert( p.total_votes >= 0, "something bad happened" );
            });
         } else {
            datxos_assert( !pd.second.second /* not from new set */, "producer is not registered" ); //data corruption
         }
      }

      _voters.modify( voter, 0, [&]( auto& av ) {
         av.last_vote_weight = new_vote_weight;
         av.producers = producers;
         av.proxy     = proxy;
      });
   }

   /**
    *  An account marked as a proxy can vote with the weight of other accounts which
    *  have selected it as a proxy. Other accounts must refresh their voteproducer to
    *  update the proxy's weight.
    *
    *  @param isproxy - true if proxy wishes to vote on behalf of others, false otherwise
    *  @pre proxy must have something staked (existing row in voters table)
    *  @pre new state must be different than current state
    */
   void system_contract::regproxy( const account_name proxy, bool isproxy ) {
      require_auth( proxy );

      auto pitr = _voters.find(proxy);
      if ( pitr != _voters.end() ) {
         datxos_assert( isproxy != pitr->is_proxy, "action has no effect" );
         datxos_assert( !isproxy || !pitr->proxy, "account that uses a proxy is not allowed to become a proxy" );
         _voters.modify( pitr, 0, [&]( auto& p ) {
               p.is_proxy = isproxy;
            });
         propagate_weight_change( *pitr );
      } else {
         _voters.emplace( proxy, [&]( auto& p ) {
               p.owner  = proxy;
               p.is_proxy = isproxy;
            });
      }
   }

   void system_contract::propagate_weight_change( const voter_info& voter ) {
      datxos_assert( voter.proxy == 0 || !voter.is_proxy, "account registered as a proxy is not allowed to use a proxy" );
      double new_weight = stake2vote( voter.staked );
      if ( voter.is_proxy ) {
         new_weight += voter.proxied_vote_weight;
      }

      /// don't propagate small changes (1 ~= epsilon)
      if ( fabs( new_weight - voter.last_vote_weight ) > 1 )  {
         if ( voter.proxy ) {
            auto& proxy = _voters.get( voter.proxy, "proxy not found" ); //data corruption
            _voters.modify( proxy, 0, [&]( auto& p ) {
                  p.proxied_vote_weight += new_weight - voter.last_vote_weight;
               }
            );
            propagate_weight_change( proxy );
         } else {
            auto delta = new_weight - voter.last_vote_weight;
            for ( auto acnt : voter.producers ) {
               auto& pitr = _producers.get( acnt, "producer not found" ); //data corruption
               _producers.modify( pitr, 0, [&]( auto& p ) {
                     p.total_votes += delta;
                     _gstate.total_producer_vote_weight += delta;
               });
            }
         }
      }
      _voters.modify( voter, 0, [&]( auto& v ) {
            v.last_vote_weight = new_weight;
         }
      );
   }

} /// namespace datxossystem
