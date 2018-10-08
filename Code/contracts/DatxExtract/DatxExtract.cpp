#include "DatxExtract.hpp"
#include <DatxioLib/multi_index.hpp>
#include <DatxioLib/chain.h>

namespace datxio
{
    /// @abi action
    void extract::recordtrx(transaction_id_type trxid, account_name handler)
    {   
        require_auth(handler);
        
        account_name producers[21]; 
        uint32_t bytes_populated = get_active_producers(producers, sizeof(account_name)*21); 
        bool Isproducer = false; 
        for (int i = 0; i < sizeof(producers)/sizeof(account_name) ;i++){ 
             if(producers[i] == handler) 
             Isproducer = true; 
        } 
        datxio_assert(Isproducer, "this func can only be called by producers");

        records trans_table(_self,_self);
        auto idx = trans_table.template get_index<N(fixed_key)>();
        auto itr = idx.find( get_fixed_key(trxid) );
        datxio_assert(itr==idx.end(), "trxid already exists");

        successtrxs success_table(_self,_self);
        auto idx2 = success_table.template get_index<N(fixed_key)>();
        auto itr2 = idx2.find(get_fixed_key(trxid) );
        datxio_assert(itr2==idx2.end(), "trxid already success");

        expirations expire_table(_self,_self);
        auto idx3 = expire_table.template get_index<N(fixed_key)>();
        auto itr3 = idx3.find(get_fixed_key(trxid) );
        if(itr3 != idx3.end()){
            expire_table.erase(*itr3);
        }

        trans_table.emplace(_self, [&](auto &s) {
            s.id= trans_table.available_primary_key();
            s.trxid= trxid;
            s.start_time = current_time();
            s.handler = handler;
        });
    }


    /// @abi action
    void extract::setverifiers(vector<account_name> accounts){
        require_auth(_self);

        verifiers veri_table(_self,_self);
        for(auto itr = veri_table.begin(); itr != veri_table.end();) {
            itr = veri_table.erase(itr);
        }

        for( const auto& a : accounts ){
            veri_table.emplace(_self, [&](auto &s) {
                s.account = a;
            });
        }
    }

    /// @abi action
    void extract::setdoing(transaction_id_type trxid, account_name handler,account_name verifier){
        require_auth(verifier);
        verifiers veri_table(_self,_self);
        auto vidx = veri_table.find(verifier);
        datxio_assert(vidx != veri_table.end(),"verifier is not exists");

        records trans_table(_self,_self);
        auto idx = trans_table.template get_index<N(fixed_key)>();
        auto itr = idx.find(get_fixed_key(trxid) );
        datxio_assert(itr != idx.end(), "trxid not exists");
        datxio_assert(itr->handler == handler, "trxid not this handler");
        
        auto itr2 = std::find( itr->verifiers.begin(), itr->verifiers.end(), verifier );
        print("verifier name",N(*itr2),"\n");
        datxio_assert( itr2 == itr->verifiers.end(), "this verifier has comfirmed this trxid" );

        trans_table.modify(*itr, get_self(), [&](auto& p)
                                              {
                                                p.verifiers.push_back(verifier);
                                              });
        if(itr->verifiers.size() < 15){
            trans_table.modify(*itr, get_self(), [&](auto& p)
                                              {
                                                p.countdown_time = current_time();
                                              });
        }
    }
    
    /// @abi action
    void extract::setsuccess(transaction_id_type trxid ,account_name verifier){
        
        require_auth(verifier);
        verifiers veri_table(_self,_self);
        auto& check = veri_table.get( verifier, "verifier not found" );

        records trans_table(_self,_self);
        auto idx = trans_table.template get_index<N(fixed_key)>();
        auto itr = idx.find(get_fixed_key(trxid) );
        datxio_assert(itr == idx.end(), "trxid not in doing records");

        trans_table.modify(*itr, get_self(), [&](auto& p)
                                              {
                                                p.successverifiers.push_back(verifier);
                                              });
        
        if (itr -> successverifiers.size() >= 15) {
            successtrxs success_table(_self,_self);
            success_table.emplace(_self, [&](auto &s) {
                s.id= success_table.available_primary_key();
                s.trxid= trxid;
                s.handler = itr -> handler;
                s.timestamp = current_time();
            });
            trans_table.erase(*itr);
        }
        
    }

    void extract::expiretrx(){
        records trans_table(_self,_self);
        auto idx = trans_table.template get_index<N(start_time)>();
        for ( auto it = idx.cbegin(); it != idx.cend();) {
            if(current_time() - it->start_time <  5*60*1000000){
                break;
            }
            if(it->verifiers.size() < 15){
                expirations expire_table(_self,_self);
                expire_table.emplace(_self, [&](auto &s) {
                    s.id= expire_table.available_primary_key();
                    s.trxid= it->trxid;
                    s.timestamp = current_time();
                    s.handler = it->handler;
                });
                trans_table.erase(*it);
            }
            else{
                 ++it; 
            }
        }
    }

    void rollbacktrx(){
        
    }

} // namespace Datxio

DATXIO_ABI( datxio::extract, (recordtrx)(setverifiers)(setdoing)(setsuccess))