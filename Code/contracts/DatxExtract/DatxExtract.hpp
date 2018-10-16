#pragma once
#include <DatxioLib/DatxioLib.hpp>
#include <DatxioLib/asset.hpp>

#include <string>

namespace datxio
{
using std::string;
class extract : public contract
{
  private:
    // void expiretrx();
    void rollbacktrx();
  public:
    extract(account_name self) : contract(self) {}

    /// @abi action
    void recordtrx(transaction_id_type trxid, account_name producer,string category);

    /// @abi action
    void setdoing(transaction_id_type trxid, account_name producer,account_name verifier);
    
    /// @abi action
    void setsuccess(transaction_id_type trxid, account_name producer);
    
    /// @abi action
    void setverifiers(vector<account_name> verifiers);

    /// @abi action
    void updateexpire();

    static key256 get_fixed_key(const checksum256& trxid) {
        const uint64_t *p64 = reinterpret_cast<const uint64_t *>(&trxid);
        return key256::make_from_word_sequence<uint64_t>(p64[0], p64[1], p64[2], p64[3]);
    };

  public:
    /// @abi table
    struct record
    {
        uint64_t            id; 
        transaction_id_type trxid;
        uint32_t                start_time;
        vector<account_name>      verifiers;
        uint32_t                countdown_time;
        vector<account_name>      successconfirm;
        account_name        producer;
        string              category;

        uint64_t primary_key() const { return id; }
        key256 by_fixed_key() const {return get_fixed_key(trxid);}
        uint64_t by_start_time() const {return start_time;}

        DATXLIB_SERIALIZE(record, (id)(trxid)(start_time)(verifiers)(countdown_time)(successconfirm)(producer)(category))
    };

    
    typedef multi_index<N(record), record,indexed_by<N(fixed_key), const_mem_fun<record, key256, &record::by_fixed_key>>
        ,indexed_by<N(start_time),const_mem_fun<record,uint64_t,&record::by_start_time>>> records;
    /// @abi table
    struct success
    {
        uint64_t            id;
        transaction_id_type trxid;
        account_name        producer;
        uint32_t                timestamp;
        string              category;            

        uint64_t primary_key() const { return id; }
        key256 by_fixed_key() const {return get_fixed_key(trxid);}

        DATXLIB_SERIALIZE(success, (id)(trxid)(producer)(timestamp)(category))
    };
    typedef datxio::multi_index<N(success), success,indexed_by<N(fixed_key), const_mem_fun<success, key256, &success::by_fixed_key>>> successtrxs;

     /// @abi table
    struct expiration
    {
        uint64_t            id; 
        transaction_id_type trxid;
        account_name        producer;
        uint32_t            timestamp;
        string              category;

        uint64_t primary_key() const { return id; }
        key256 by_fixed_key() const {return get_fixed_key(trxid);}

        DATXLIB_SERIALIZE(expiration, (id)(trxid)(producer)(timestamp)(category))
    };
    typedef datxio::multi_index<N(expiration), expiration,indexed_by<N(fixed_key), const_mem_fun<expiration, key256, &expiration::by_fixed_key>>> expirations;


    /// @abi table
    struct verifier
    {
        account_name account;

        uint64_t primary_key() const { return account; }
        DATXLIB_SERIALIZE(verifier, (account))
    };
    typedef datxio::multi_index<N(verifier), verifier> verifiers;

    /// @abi table
    struct transrecord
    {
        uint64_t            id; 
        transaction_id_type trxid;
        string              category;
        account_name        account;
        asset               quantity;
        string              memo;

        uint64_t primary_key() const { return id; }
        key256 by_fixed_key() const {return get_fixed_key(trxid);}

        DATXLIB_SERIALIZE(transrecord, (id)(trxid)(category)(account)(quantity)(memo))
    };
    typedef datxio::multi_index<N(transrecord), transrecord,indexed_by<N(fixed_key), const_mem_fun<transrecord, key256, &transrecord::by_fixed_key>>> transrecords;


};

} // namespace Datxio
