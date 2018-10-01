#include <DatxioLib/types.hpp>
#include <DatxioLib/currency.hpp>
#include <boost/container/flat_map.hpp>
#include <cmath>
#include <DatxExchange/market_state.hpp>

namespace datxio {

   
   class exchange {
      private:
         account_name      _this_contract;
         currency          _excurrencies;
         exchange_accounts _accounts;

      public:
         exchange( account_name self )
         :_this_contract(self),
          _excurrencies(self),
          _accounts(self)
         {}

         void createx( account_name    creator,
                       asset           initial_supply,
                       uint32_t        fee,
                       extended_asset  base_deposit,
                       extended_asset  quote_deposit
                     );

         void deposit( account_name from, extended_asset quantity );
         void withdraw( account_name  from, extended_asset quantity );
         void lend( account_name lender, symbol_type market, extended_asset quantity );

         void unlend(
            account_name     lender,
            symbol_type      market,
            double           interest_shares,
            extended_symbol  interest_symbol
         );

         struct covermargin {
            account_name     borrower;
            symbol_type      market;
            extended_asset   cover_amount;
         };

         struct upmargin {
            account_name     borrower;
            symbol_type      market;
            extended_asset   delta_borrow;
            extended_asset   delta_collateral;
         };

         struct trade {
            account_name    seller;
            symbol_type     market;
            extended_asset  sell;
            extended_asset  min_receive;
            uint32_t        expire = 0;
            uint8_t         fill_or_kill = true;
         };

         void on( const trade& t    );
         void on( const upmargin& b );
         void on( const covermargin& b );
         void on( const currency::transfer& t, account_name code );

         void apply( account_name contract, account_name act );
   };
} // namespace datxio
