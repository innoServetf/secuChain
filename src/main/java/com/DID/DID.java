package com.DID;

import lombok.Data;

@Data
public class DID {
    public String DID;
    public KeyPairs authKey;
    public KeyPairs recyKey;
    public DIDDocument document;

}
