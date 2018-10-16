#include "DatxRecharge.hpp"
#include <DatxioLib/multi_index.hpp>
#include <fc/string.hpp>
#include <cstring>
#include <DatxioLib/crypto.h>
#include <fc/crypto/hex.hpp>
#include <DatxioLib/chain.h>
namespace datxio
{
/// @abi action
void recharge::recorduser(account_name datxaddress,string address,account_name bpname)
{
    // require_auth(bpname);
    // account_name producers[21]; 
    // uint32_t bytes_populated = get_active_producers(producers, sizeof(account_name)*21); 
    //     bool Isproducer = false; 
    //     for (int i = 0; i < sizeof(producers)/sizeof(account_name) ;i++){ 
    //          if(producers[i] == bpname) 
    //          Isproducer = true; 
    //     } 
    //     datxio_assert(Isproducer, "this func can only be called by producers");

    // define the table 
        users users_table(_self,_self);

        char *c=(char*)address.c_str();
        checksum256 calc_hash;
        sha256( c, sizeof(c), &calc_hash );

        auto idx = users_table.template get_index<N(address)>();
        auto itr = idx.find( get_hash(calc_hash) );
        datxio_assert(itr==idx.end(), "This address is already exists");
        users_table.emplace(_self, [&](auto &u) {
            u.userid= users_table.available_primary_key();
            u.datxaddress=datxaddress;
            u.address=calc_hash;
        });    
}
/// @abi action
void recharge::charge(account_name bpname,string hash,string from,string to,int64_t blocknum,string quantity,string category,string memo)
 {
//         require_auth(bpname);
//     //check the caller is producer
//         account_name producers[21]; 
//         uint32_t bytes_populated = get_active_producers(producers, sizeof(account_name)*21); 
//         bool Isproducer = false; 
//         for (int i = 0; i < sizeof(producers)/sizeof(account_name) ;i++){ 
//              if(producers[i] == bpname) 
//              Isproducer = true; 
//         } 
//         datxio_assert(Isproducer, "this func can only be called by producers");
        
    //get all expired transactions and remove to expirdrecord
        this->expired_trx();
    //check the countrecord whether exist
        countrecords count_table(_self,_self);
        perfections success_table(_self,_self);
        expirations expired_table(_self,_self);
        records records_table(_self,_self);

    // caculate hash for trx data

        string str=hash+from+to+std::to_string(blocknum)+quantity+category+memo;
        char *data=(char*)str.c_str();
        checksum256 calc_hash;
        sha256( data, sizeof(data), &calc_hash );

    //check if the hash alerday success
        auto idx2 = success_table.template get_index<N(data)>();
        auto itr2 = idx2.find(get_hash(calc_hash) );
        datxio_assert(itr2==idx2.end(), "This transaction is already success");

    //check the data hash whether equal
        auto idx = records_table.template get_index<N(data)>();
        auto itr = idx.find(get_hash(calc_hash));
        if(itr!=idx.end()){
            //check producer whether more than 15
            if(sizeof(itr->producers)<15)
            {
                auto itr3 = std::find( itr->producers.begin(), itr->producers.end(), bpname );
                datxio_assert( itr3 != itr->producers.end(), "This producer already initiated a request for this transaction!");    
                records_table.modify(*itr, get_self(), [&](auto& p)
                                                {
                                                    p.producers.push_back(bpname);
                                                });
            }else
            {
            //transanfer to datx account 
            // you should add DatxRecharge permission first, and then can send action DatxToken contract.

                char *c=(char*)from.c_str();
                checksum256 address;
                sha256( c, sizeof(c), &address);
                asset quant{fc::to_int64(quantity),S(4,category)};
                if(category=="EOS"){
                    action(
                        permission_level{N(datxio.deos), N(active) },
                        N(datxio.deos), N(transfer),
                        std::make_tuple(N(datxio.deos), N(memo), quant,memo)
                    ).send();
                }else{
                    users user_table(_self,_self);
                    auto index = user_table.template get_index<N(address)>();
                    auto addr = index.find( get_hash(address));
                    datxio_assert(addr==index.end(), "This address isn't exists");
                    if(category=="BTC")
                    {
                        action(
                            permission_level{N(datxio.dbtc), N(active) },
                            N(datxio.dbtc), N(transfer),
                            std::make_tuple(N(datxio.dbtc), N(addr->datxaddress), quant,memo)
                        ).send();
                    }else
                    {
                        action(
                            permission_level{N(datxio.deth), N(active) },
                            N(datxio.deth), N(transfer),
                            std::make_tuple(N(datxio.deth), N(addr->datxaddress), quant,memo)
                        ).send();
                    }
                }
                //after transfer to datxaddress remove info to success
                success_table.emplace(_self, [&](auto &s) 
                {
                    s.id= records_table.available_primary_key();
                    s.trxid= hash;
                    s.from = from;
                    s.to = to;
                    s.blocknum= std::to_string(blocknum);
                    s.quantity=quantity;
                    s.category=category;
                    s.memo = memo;
                    s.data=calc_hash; 
                });   
                records_table.erase(*itr);//remove the transaction in the record_table
                auto in = expired_table.template get_index<N(data)>();
                auto r = in.find(get_hash(calc_hash) );
                if(r!=in.end())
                {
                    expired_table.erase(*r);
                }
            }
            
        }else
        {
            // if data hash isn't exist then record it in record table 
            vector<account_name> producers_name;
            producers_name.push_back(bpname);
            records_table.emplace(_self, [&](auto &s) {
                s.id= records_table.available_primary_key();
                s.trxid= hash;
                s.from = from;
                s.to = to;
                s.blocknum= blocknum;
                s.quantity=quantity;
                s.category=category;
                s.memo = memo;
                s.data=calc_hash; 
                s.start_time=current_time(); //seconds
                s.producers=producers_name;
            });       
        }
}

//get the expired transaction from record
void recharge::expired_trx()
{
    records trans_table(_self,_self);
    int count = 0;
    for ( auto it = trans_table.cbegin(); it != trans_table.cend() && count<100;)
    { 
        uint64_t subtime = now() - it->start_time;
        if(it->producers.size() < 15&&subtime > 5*63)
        {
            countrecords count_table(_self,_self);
            auto id =count_table.template get_index<N(data)>();
            auto i = id.find( get_hash(it->data) );
  
                if(i!=id.end()&&i->count <3)
                {
                    count_table.modify(*i, get_self(), [&](auto& c)
                    {
                        c.count=i->count+1;
                    });

                    ++it;
                }else if (i==id.end())
                {
                    count_table.emplace(_self, [&](auto &co){
                        co.id=count_table.available_primary_key();
                        co.data=it->data;
                        co.count=1;
                    });

                    expirations expire_table(_self,_self);
                    expire_table.emplace(_self, [&](auto &e) {
                        e.id= expire_table.available_primary_key();
                        e.trxid=it->trxid;
                        e.from=it->from;
                        e.blocknum=it->blocknum;
                        e.quantity=it->quantity;
                        e.category=it->category;
                        e.memo=it->memo;
                        e.data=it->data;        
                    }); 
                    it = trans_table.erase(it);
                    ++count;
                }else{
                    ++it;
                }
        }else
        {
            ++it;
        }
    }
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

DATXIO_ABI_EX(datxio::recharge,(recorduser)(charge))
