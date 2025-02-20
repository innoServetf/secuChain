package com.DID;

import com.google.gson.Gson;

import java.security.*;
import java.security.spec.ECGenParameterSpec;
import java.security.spec.PKCS8EncodedKeySpec;
import java.time.Instant;
import java.util.Base64;
import java.util.UUID;

import static com.google.common.hash.Hashing.sha256;


public class DIDClient {

    private static final Gson gson = new Gson();


    public String generateDID() {
        String specificId = UUID.randomUUID().toString();
        return "did:" + "fabric:" + specificId;
    }



    //密鑰生成
    public KeyPair generateKeyPair2() throws NoSuchAlgorithmException, InvalidAlgorithmParameterException {
        KeyPairGenerator keyGen = KeyPairGenerator.getInstance("EC");
        ECGenParameterSpec ecSpec = new ECGenParameterSpec("secp256k1");
        keyGen.initialize(ecSpec, new SecureRandom());
        return keyGen.generateKeyPair();
    }

    public KeyPair generateKeyPair() throws NoSuchAlgorithmException {
        KeyPairGenerator keyGen = KeyPairGenerator.getInstance("RSA");
        keyGen.initialize(2048, new SecureRandom());
        return keyGen.generateKeyPair();
    }

    public String encodePublicKey(KeyPair keyPair) {
        return Base64.getEncoder().encodeToString(keyPair.getPublic().getEncoded());
    }

    public String encodePrivateKey(KeyPair keyPair) {
        return Base64.getEncoder().encodeToString(keyPair.getPrivate().getEncoded());
    }
    /**
     * 用私钥签名 DID Document
     */
    public  String signDidDocument(PrivateKey privateKey, DIDDocument didDocument) throws Exception {
        Gson gson = new Gson();
        Signature signature = Signature.getInstance("SHA256withECDSA", "BC");
        signature.initSign(privateKey);
        String documentJson =gson.toJson(didDocument);
        signature.update(documentJson.getBytes());
        byte[] signatureBytes = signature.sign();
        return Base64.getEncoder().encodeToString(signatureBytes);
    }




    //驗證公私鑰是否為空
    public boolean verifyRecoveryKey(String publicKey,String privateKey){
        return publicKey!=null&&privateKey!=null;
    }



    public DID createDID() throws Exception {

        //1
        DID did = new DID();
        did.setDID(generateDID());
        //2主
        KeyPair keypair = generateKeyPair();
        KeyPairs keyPairs = new KeyPairs();
        keyPairs.setType("RSA");
        PrivateKey privateKey = keypair.getPrivate();
        keyPairs.setPublicKey(encodePublicKey(keypair));
        keyPairs.setPrivateKey(encodePrivateKey(keypair));
        did.setAuthKey(keyPairs);


        //3備
        KeyPair keyPair2 = generateKeyPair2();
        KeyPairs keyPairs2 = new KeyPairs();
        keyPairs2.setType("EC");
        keyPairs2.setPublicKey(encodePublicKey(keyPair2));
        keyPairs2.setPrivateKey(encodePrivateKey(keyPair2));
        did.setRecyKey(keyPairs2);


        //4
        DIDDocument didDocument = new DIDDocument();
        did.setDocument(generatorDIDDocument(did.DID,keyPairs.publicKey,keyPairs2.publicKey));

        Proof proof = new Proof();
        proof.setCreator(did.getDID());
        proof.setType("EcdsaSecp256k1Signature");
        proof.setSignature(signDidDocument(privateKey,didDocument));

        did.document.created= Instant.now().toString();
        did.document.updated =did.document.created;
        return did;


    }


    //生成DID文檔
    public DIDDocument generatorDIDDocument(String did,String publicKey,String publicKey2) throws NoSuchAlgorithmException {
        DIDDocument didDocument = new DIDDocument();
        //1
        didDocument.setDid(did);
        //2
        Authentication authentication = new Authentication();
        authentication.setType("RSA");
        authentication.setPublicKey(publicKey);
        didDocument.setAuthentication(authentication);
        //3
        Recovery recovery = new Recovery();
        recovery.setType("EC");
        recovery.setPublicKey(publicKey2);
        didDocument.setRecovery(recovery);


        return didDocument;

    }





}
