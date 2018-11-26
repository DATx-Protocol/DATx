#include "DatxExtract.hpp"
#include <DatxosLib/multi_index.hpp>
#include <DatxosLib/chain.h>

namespace datxos
{
    /// @abi action
    void extract::recordtrx(transaction_id_type trxid, account_name producer,string category)
    {   
        require_auth(producer);
        account_name producers[3]; 
        auto bytes_populated = get_active_producers(producers, sizeof(producers)); 
        bool Isproducer = false; 
        for (size_t i = 0; i < bytes_populated/ sizeof(account_name) ;i++){ 
            if(producers[i] == producer) {
                Isproducer = true; 
                break;
            } 
        } 
        datxos_assert(Isproducer, "this func can only be called by producers");

        if (category != "DBTC" && category != "DETH" && category != "DEOS"){
            datxos_assert(false, "the category is not (DBTC,DETH,DEOS).");
        }

        transrecords transrecord_table(N(datxos.dtoke),N(datxos.dtoke));
        auto idx0 = transrecord_table.template get_index<N(fixed_key)>();
        auto itr0 = idx0.find( get_fixed_key(trxid) );
        datxos_assert(itr0 != idx0.end(), "trxid not exists in token transfer");

        std::string newCategory = category;
        datxos_assert(itr0 ->quantity.symbol == string_to_symbol(4,newCategory.c_str()),"category not correct");
 
        records trans_table(_self,_self);
        auto idx = trans_table.template get_index<N(fixed_key)>();
        auto itr = idx.find( get_fixed_key(trxid) );
        datxos_assert(itr==idx.end(), "trxid already exists");

        successtrxs success_table(_self,_self);
        auto idx2 = success_table.template get_index<N(fixed_key)>();
        auto itr2 = idx2.find(get_fixed_key(trxid) );
        datxos_assert(itr2==idx2.end(), "trxid already success");

        expirations expire_table(_self,_self);
        auto idx3 = expire_table.template get_index<N(fixed_key)>();
        auto itr3 = idx3.find(get_fixed_key(trxid) );
        if(itr3 != idx3.end()){
            expire_table.erase(*itr3);
        }

        trans_table.emplace(_self, [&](auto &s) {
            s.id= trans_table.available_primary_key();
            s.trxid= trxid;
            s.start_time = now();
            s.producer = producer;
            s.category = newCategory;
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
    void extract::setdoing(transaction_id_type trxid, account_name producer,account_name verifier){
        // require_auth(verifier);
        verifiers veri_table(_self,_self);
        auto vidx = veri_table.find(verifier);
        datxos_assert(vidx != veri_table.end(),"verifier is not exists");

        records trans_table(_self,_self);
        auto idx = trans_table.template get_index<N(fixed_key)>();
        auto itr = idx.find(get_fixed_key(trxid) );
        datxos_assert(itr != idx.end(), "trxid not exists");
        datxos_assert(itr->producer == producer, "trxid not this producer");
        
        auto itr2 = std::find( itr->verifiers.begin(), itr->verifiers.end(), verifier );
        datxos_assert( itr2 == itr->verifiers.end(), "this verifier has comfirmed this trxid" );

        trans_table.modify(*itr, get_self(), [&](auto& p)
                                              {
                                                p.verifiers.push_back(verifier);
                                              });
        
        account_name producers[3];
        auto bytes_populated = get_active_producers(producers, sizeof(producers));      
        int psize = bytes_populated/sizeof(account_name)*2/3;

        if(itr->verifiers.size() > psize){
            trans_table.modify(*itr, get_self(), [&](auto& p)
                                              {
                                                p.countdown_time = now();
                                              });
        }
    }
    
    /// @abi action
    void extract::setsuccess(transaction_id_type trxid ,account_name producer){
        require_auth(producer);

        account_name producers[3]; 
        auto bytes_populated = get_active_producers(producers, sizeof(producers)); 
        int psize = bytes_populated/sizeof(account_name)*2/3;
        bool Isproducer = false; 
        for (size_t i = 0; i < bytes_populated/sizeof(account_name) ;i++){ 
            if(producers[i] == producer) {
                Isproducer = true; 
                break;
            }
        } 
        datxos_assert(Isproducer, "this func can only be called by producers");

        records trans_table(_self,_self);
        auto idx = trans_table.template get_index<N(fixed_key)>();
        auto itr = idx.find(get_fixed_key(trxid) );
        datxos_assert(itr != idx.end(), "trxid not in doing records");

        auto itr2 = std::find( itr->successconfirm.begin(), itr->successconfirm.end(), producer );
        datxos_assert( itr2 == itr->successconfirm.end(), "setsuccess this producer has comfirmed this trxid" );

        trans_table.modify(*itr, get_self(), [&](auto& p)
                                              {
                                                p.successconfirm.push_back(producer);
                                              });
        
        if (itr -> successconfirm.size() > psize) {
            successtrxs success_table(_self,_self);
            success_table.emplace(_self, [&](auto &s) {
                s.id= success_table.available_primary_key();
                s.trxid= trxid;
                s.producer = itr -> producer;
                s.timestamp = now();
                s.category = itr->category;
            });
            trans_table.erase(*itr);
        }
        
    }

    /// @abi action
    void extract::updateexpire(){
         //this->rollbacktrx();

        account_name producers[3];
        uint32_t bytes_populated = get_active_producers(producers, sizeof(producers));      
        int psize = bytes_populated/sizeof(account_name)*2/3;
        records trans_table(_self,_self);
        auto idx = trans_table.template get_index<N(start_time)>();
        int count = 0;
        for ( auto it = idx.cbegin(); it != idx.cend() && count < 5;) {
            uint64_t subtime = now() - it->start_time;
            if(subtime <= 5 * 60){
                break;
            }
            if(it->verifiers.size() < psize){
                expirations expire_table(_self,_self);
                expire_table.emplace(_self, [&](auto &s) {
                    s.id= expire_table.available_primary_key();
                    s.trxid= it->trxid;
                    s.timestamp = now();
                    s.producer = it->producer;
                    s.category = it->category;
                });

                it = idx.erase(it);

                ++count;

            }else{
                ++it;
            }
        } 
    }
/// @abi action
    void extract::rollbacktrx(){
        records trans_table(_self,_self);
        successtrxs success_table(_self,_self);
        int count = 0; 
        for ( auto it = trans_table.cbegin(); it != trans_table.cend()&&count < 5;) { 
            auto trxid = it->trxid;

            //get transaction details and rollback
            if (it->countdown_time == 0) {
                continue;
            }
            uint64_t subtime = now() - it->countdown_time;

            if ((it->category == "DETH" && subtime > 30*60) || (it->category == "DEOS" && subtime > 5*60)){
                        
                transrecords transrecord_table(N(datxos.dtoke),N(datxos.dtoke));
                auto idx0 = transrecord_table.template get_index<N(fixed_key)>();
                auto itr0 = idx0.find( get_fixed_key(trxid) );

                datxos_assert(itr0 != idx0.end(), "trxid not exists in token transfer");

                // pay the cost
                if(it->category == "DETH"){
                    action(permission_level{N(datxos.deth), N(active) },
                        N(datxos.dtoke), N(transfer),
                        std::make_tuple(N(datxos.deth),itr0 -> account, itr0 -> quantity, std::string("rollback"))
                    ).send();
                }
                else if (it->category == "DEOS"){
                    action(permission_level{N(datxos.deos), N(active) },
                        N(datxos.dtoke), N(transfer),
                        std::make_tuple(N(datxos.deos),itr0 -> account, itr0 -> quantity, std::string("rollback"))
                    ).send();
                }

                success_table.emplace(_self, [&](auto &s) {
                    s.id= success_table.available_primary_key();
                    s.trxid= trxid;
                    s.producer = it -> producer;
                    s.timestamp = now();
                    s.category = it->category;
                });
                it = trans_table.erase(it);
                ++count;
            }else{
                ++it;
            }
        }
    }

} //namespace Datxos

DATXOS_ABI( datxos::extract, (recordtrx)(setverifiers)(setdoing)(setsuccess)(updateexpire)(rollbacktrx))