/**
 *  @file
 *  @copyright defined in datx/LICENSE.txt
 */
#pragma once

#include <DatxosLib/asset.hpp>
#include <DatxosLib/DatxosLib.hpp>

#include <string>

namespace datxossystem {
   class system_contract;
}

namespace datxos {

   using std::string;

   class token : public contract {
       static key256 get_fixed_key(const checksum256& trxid) {
        const uint64_t *p64 = reinterpret_cast<const uint64_t *>(&trxid);
        return key256::make_from_word_sequence<uint64_t>(p64[0], p64[1], p64[2], p64[3]);
    };

      public:
         token( account_name self ):contract(self){}

         void create( account_name issuer,
                      asset        maximum_supply);

         void issue( account_name to, asset quantity, string memo );

         void transfer( account_name from,
                        account_name to,
                        asset        quantity,
                        string       memo );
      
          void extract( account_name from,
                        account_name to,
                        asset        quantity,
                        string       memo );

         inline asset get_supply( symbol_name sym )const;
         
         inline asset get_balance( account_name owner, symbol_name sym )const;

      private:
         struct account {
            asset    balance;

            uint64_t primary_key()const { return balance.symbol.name(); }
         };

         struct currency_stats {
            asset          supply;
            asset          max_supply;
            account_name   issuer;

            uint64_t primary_key()const { return supply.symbol.name(); }
         };

         typedef datxos::multi_index<N(accounts), account> accounts;
         typedef datxos::multi_index<N(stat), currency_stats> stats;

         void sub_balance( account_name owner, asset value );
         void add_balance( account_name owner, asset value, account_name ram_payer );

      public:
         struct transfer_args {
            account_name  from;
            account_name  to;
            asset         quantity;
            string        memo;
         };

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
    typedef datxos::multi_index<N(transrecord), transrecord,indexed_by<N(fixed_key), const_mem_fun<transrecord, key256, &transrecord::by_fixed_key>>> transrecords;

   };

   asset token::get_supply( symbol_name sym )const
   {
      stats statstable( _self, sym );
      const auto& st = statstable.get( sym );
      return st.supply;
   }

   asset token::get_balance( account_name owner, symbol_name sym )const
   {
      accounts accountstable( _self, owner );
      const auto& ac = accountstable.get( sym );
      return ac.balance;
   }

} /// namespace datxos
