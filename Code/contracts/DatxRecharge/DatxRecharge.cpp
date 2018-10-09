#include "DatxRecharge.hpp"
#include <DatxioLib/multi_index.hpp>
#include <fc/time.hpp>
#include <fc/string.hpp>
#include <cstring>
#include <DatxioLib/crypto.h>
namespace datxio
{
/// @abi action
void recharge:recorduser(account_name datxaddress,string address,account_namebpname)
{
     require_auth(bpname);
        account_name producers[21]; 
        uint32_t bytes_populated = get_active_producers(producers, sizeof(account_name)*21); 
        bool Isproducer = false; 
        for (int i = 0; i < sizeof(producers)/sizeof(account_name) ;i++){ 
             if(producers[i] == bpname) 
             Isproducer = true; 
        } 
        datxio_assert(Isproducer, "this func can only be called by producers");

    // define the table 
        users users_table(_self,_self);
         auto idx = records_table.template get_index<N(address)>();
         auto itr = idx.find( user::get_hash(address) );
        datxio_assert(itr==idx.end(), "This address is already exists");
        users_table.emplace(_self, [&](auto &u) {
            u.userid= users_table.available_primary_key();
            u.datxaddress=datxaddress;
            u.address=address;
        });
}

void recharge::charge(account_name bpname,transaction_id_type hash,string from,string to,int64_t blocknum,string quantity,string category,string memo)
{
        require_auth(bpname);
        string txHash = to_hex((const char*)(&hash), sizeof(hash));
    //check the caller is producer
        account_name producers[21]; 
        uint32_t bytes_populated = get_active_producers(producers, sizeof(account_name)*21); 
        bool Isproducer = false; 
        for (int i = 0; i < sizeof(producers)/sizeof(account_name) ;i++){ 
             if(producers[i] == bpname) 
             Isproducer = true; 
        } 
        datxio_assert(Isproducer, "this func can only be called by producers");
        
       
    //get all expired transactions and remove to expirdrecord
        expired_trx();
    //check the countrecord whether exist
        
        countrecords count_table(_self,_self);
        successrecords success_table(_self,_self);
        expiredrecords expired_table(_self,_self);
         records records_table(_self,_self);
        
        

//(3)check the hash whether exits
    // caculate hash for trx data
        string str=txhash+from+to+to_string(blocknum)+quantity+category+memo;
        char *data=(char*)str.c_str();
        checksum calc_hash;
        sha256( data, length, &calc_hash );
    //check if the hash alerday success
        auto idx2 = success_table.template get_index<N(data)>();
        auto itr2 = idx2.find(record::by_trxid(calc_hash) );
        datxio_assert(itr2==idx2.end(), "This transaction is already success");

    //check the data hash whether equal
        auto idx = records_table.template get_index<N(data)>();
        auto itr = idx.find(record::by_data(calc_hash));
        if(itr==idx.end()){


    //check producer whether more than 15
        bool bpexist=false;
        for( auto it = itr->producers.begin(); it !=itr->producers.end(); ++it )
        {
            if(itr->producers[i] == bpname)
            {
                bpexist=true;
            }
            datxio_assert(bpexit, "This producer already initiated a request for this transaction!");
        }
        records_table.modify(itr,0, [&](auto &s) 
        {
               s.producers=s.producers.push_back(bpname)

        }

//(6)transanfer to datx account
        
         // you should add DatxRecharge permission first, and then can send action DatxToken contract.
       
        if(sizeof(itr.producers>=14))
        {

            assert quant{to_int64(quantity),S(4,category)}
            if(category=="BTC"||category=="ETH"){
                users user_table(_self,_self);
                auto index = users_table.template get_index<N(address)>();
                auto it = idx.find( user::get_hash(address) );
                bool address=false;
                if(it==index.end())
                {
                    address=true;
                }
                datxio_assert(address, "This address isn't already exists");
                if(address==true)
                {
                    
                    if(category=="BTC")
                    {
                        action(
                            permission_level{N(datxio.dbtc), N(active) },
                            N(datxio.dbtc), N(transfer),
                            std::make_tuple(N(datxio.dbtc), it, quant,memo)
                        ).send();
                    }else{
                         action(
                            permission_level{N(datxio.deth), N(active) },
                            N(datxio.deth), N(transfer),
                            std::make_tuple(N(datxio.deth), it, quant,memo)
                        ).send();
                    }
                }
            }else{
                action(
                            permission_level{N(datxio.deos), N(active) },
                            N(datxio.deos), N(transfer),
                            std::make_tuple(N(datxio.deos), N(memo), quant,memo)
                        ).send();
            }
        }

        success_table.emplace(_self, [&](auto &s) 
        {
            s.id= records_table.available_primary_key();
            s.trxid= hash;
            s.from = from;
            s.to = to;
            s.blocknum= blocknum;
            s.quantity=quantity;
            s.category=category;
            s.memo = memo;
            s.data=calc_hash; 
        }
        
        
        records_table.erase(itr);//remove the transaction in the record_table
        auto in = expired_table.template get_index<N(data)>();
        auto i = idx.find( expiredrecord::by_data(calc_hash) );
        if(in==i)
        {
            expired_table.erase(i);
        }
        }else{
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
                s.time=time_point_sec(now()); //seconds
                memcpy(s.producers,producers_name,sizeof(producers_name));
            });

        }

}

//get the expired transaction from record
void recharge::expired_trx()
{
    records trans_table(_self,_self);
    auto idx = trans_table.template get_index<N(start_time)>();
    for ( auto it = idx.cbegin(); it != idx.cend();)
    {
        if(current_time() - it->start_time < 5*60*1000000)
        {
            break;
        }
        if(it->producers.size() < 15)
        {
            countrecords count_table(_self,_self);
            auto id =count_table.template get_index<N(data)>();
            auto i = id.find( expiredrecord::by_data(it->data) );
            if(i==id.end()&&i->count >=3)
            {  
                break;
            }else
            {
                if(i==id.end()&&i->count <3)
                {
                    count_table.modify(itr, get_self(), [&](auto& c)
                    {
                        c.count=i->count+1;
                    });
                }else
                {
                    count_table.emplace(_self, [&](auto &co){
                        co.id=count_table.available_primary_key();
                        co.data=it->data;
                        co.count=1;
                    });
                }
                expiredrecords expire_table(_self,_self);
                expire_table.emplace(_self, [&](auto &e) {
                    e.id= success_table.available_primary_key();
                    e.trxid=it->trxid;
                    e.from=it->from;
                    e.blocknum=it->blocknum;
                    e.quantity=it->quantity;
                    e.category=it->category;
                    e.memo=it->memo;
                    e.data=it->data;        
                });
            }
            trans_table.erase(it);
        }
       
        else
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
