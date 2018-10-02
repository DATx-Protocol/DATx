#include "DatxRecharge.hpp"
#include <DatxioLib/multi_index.hpp>
namespace datxio
{
/// @abi action

void recharge::charge(transaction_id_type hash,string from,string to,int64_t blocknum,string quantity,string category,string memo)
{

   // define the table
        require_auth(_self);
        records records_table(_self,_self);
         auto idx = records_table.template get_index<N(hash)>();
         auto itr = idx.find( record::get_hash(hash) );
        datxio_assert(itr==idx.end(), "hash for transaction already exists");

        records_table.emplace(_self, [&](auto &s) {
            s.trxid= records_table.available_primary_key();
            s.hash= hash;
            s.from = from;
            s.to = to;
            s.blocknum= blocknum;
            s.quantity=quantity;
            s.category=category;
            s.memo = memo;
        });

}

void recharge::transtoken(transaction_id_type hash,account_name from,account_name to,asset quantity,string memo)
{
     
    require_auth(_self);
        transactions trans_table(_self,_self);
        auto idx = trans_table.template get_index<N(hash)>();
        auto itr = idx.find( transaction::get_hash(hash) );
        datxio_assert(itr==idx.end(), "hash for transaction already exists");
        trans_table.emplace(_self, [&](auto &s) {
            s.trxid= trans_table.available_primary_key();
            s.hash= hash;
            s.from = from;
            s.to = to;
            s.memo = memo;
        });
        // you should add DatxRecharge permission first, and then can send action DatxToken contract.
        action(
            permission_level{from, N(active) },
           from, N(transfer),
            std::make_tuple(from, to, quantity,memo)
        ).send();
}

} // namespace Datxio


// modify form DATXIO_ABI
#define DATXIO_ABI_EX( TYPE, MEMBERS ) \
extern "C" { \
   void apply( uint64_t receiver, uint64_t code, uint64_t action ) { \
      auto self = receiver; \
      if( action == N(onerror)) { \
         /* onerror is only valid if it is for the "Datxio" code account and authorized by "Datxio"'s "active permission */ \
         datxio_assert(code == N(Datxio), "onerror action's are only valid from the \"Datxio\" system account"); \
      } \
      if( code == self || code == N(Datxio.token) || action == N(onerror) ) { \
         TYPE thiscontract( self ); \
         switch( action ) { \
            DATXIO_API( TYPE, MEMBERS ) \
         } \
         /* does not allow destructor of thiscontract to run: Datxio_exit(0); */ \
      } \
   } \
}

DATXIO_ABI_EX(datxio::recharge,(charge)(transtoken))
