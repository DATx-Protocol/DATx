/**
 *  @file
 *  @copyright defined in datx/LICENSE.txt
 */
#include <datxos/template_accessory/template_accessory.hpp>

namespace datxos {
   static appbase::abstract_accessory& _template_accessory = app().register_accessory<template_accessory>();

class template_accessory_impl {
   public:
};

template_accessory::template_accessory():my(new template_accessory_impl()){}
template_accessory::~template_accessory(){}

void template_accessory::set_program_options(options_description&, options_description& cfg) {
   cfg.add_options()
         ("option-name", bpo::value<string>()->default_value("default value"),
          "Option Description")
         ;
}

void template_accessory::accessory_initialize(const variables_map& options) {
   try {
      if( options.count( "option-name" )) {
         // Handle the option
      }
   }
   FC_LOG_AND_RETHROW()
}

void template_accessory::accessory_startup() {
   // Make the magic happen
}

void template_accessory::accessory_shutdown() {
   // OK, that's enough magic
}

}
