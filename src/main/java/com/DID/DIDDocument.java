package com.DID;

import lombok.Data;

@Data
public  class DIDDocument {

    public String did;
    public Authentication authentication;
    public Recovery recovery;
    public Proof proof;


    public String created;
    public String updated;




}