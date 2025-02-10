package com.DID;
import com.google.gson.Gson;
import org.hyperledger.fabric.contract.Context;

import org.hyperledger.fabric.contract.ContractInterface;
import org.hyperledger.fabric.contract.annotation.Contract;
import org.hyperledger.fabric.contract.annotation.Transaction;

import java.security.*;
import java.time.Instant;

///
@Contract
public class DIDContract implements ContractInterface {

    private final Gson gson=new Gson();
    private final DIDClient didClient;

    public DIDContract(DIDClient didClient) {
        this.didClient = didClient;
    }


    @Transaction(intent = Transaction.TYPE.SUBMIT)
    public String registerDID(Context ctx) throws NoSuchAlgorithmException {
        String didKey =didClient.generateDID();
        byte[] existing = ctx.getStub().getState(didKey);
        if (existing != null && existing.length > 0) {
            throw new RuntimeException("DID already exists");
        }

        String publicKey= didClient.encodePublicKey(didClient.generateKeyPair());

        DIDDocument doc = new DIDDocument();
        doc.id = didKey;
        doc.publicKey = publicKey;
        doc.created=Instant.now().toString();
        doc.updated = doc.created;

        String json = gson.toJson(doc);
        ctx.getStub().putState(didKey, json.getBytes());
        return json;
    }




    @Transaction(intent = Transaction.TYPE.SUBMIT)
    public String updateDID(Context ctx, String did, String publicKey) {
        String didKey = did;
        byte[] existing = ctx.getStub().getState(didKey);
        if (existing == null || existing.length == 0) {
            throw new RuntimeException("DID does not exist");
        }

        DIDDocument doc = gson.fromJson(new String(existing), DIDDocument.class);
        doc.publicKey = publicKey;
        doc.updated = Instant.now().toString();


        String json = gson.toJson(doc);
        ctx.getStub().putState(didKey, json.getBytes());
        return json;
    }


    @Transaction(intent = Transaction.TYPE.EVALUATE)
    public String queryDID(Context ctx, String did) {
        String didKey = "did:" + "fabric:" + did;
        byte[] existing = ctx.getStub().getState(didKey);
        if (existing == null || existing.length == 0) {
            throw new RuntimeException("DID does not exist");
        }
        return new String(existing);
    }




}