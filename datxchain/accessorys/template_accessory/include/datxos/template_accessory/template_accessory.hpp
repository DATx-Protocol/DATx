/**
 *  @file
 *  @copyright defined in datx/LICENSE.txt
 */
#pragma once
#include <appbase/application.hpp>

namespace datxos {

using namespace appbase;

/**
 *  This is a template accessory, intended to serve as a starting point for making new accessorys
 */
class template_accessory : public appbase::accessory<template_accessory> {
public:
   template_accessory();
   virtual ~template_accessory();
 
   APPBASE_accessory_REQUIRES()
   virtual void set_program_options(options_description&, options_description& cfg) override;
 
   void accessory_initialize(const variables_map& options);
   void accessory_startup();
   void accessory_shutdown();

private:
   std::unique_ptr<class template_accessory_impl> my;
};

}
