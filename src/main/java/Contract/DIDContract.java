package Contract;

import com.DID.DID;
import com.DID.DIDClient;
import com.DID.DIDDocument;
import com.DID.KeyPairs;
import com.google.gson.Gson;
import org.hyperledger.fabric.contract.Context;
import org.hyperledger.fabric.contract.ContractInterface;
import org.hyperledger.fabric.contract.annotation.Contract;
import org.hyperledger.fabric.contract.annotation.Transaction;

import java.nio.charset.StandardCharsets;
import java.security.InvalidAlgorithmParameterException;
import java.security.KeyPair;
import java.security.NoSuchAlgorithmException;
import java.time.Instant;

@Contract
public class DIDContract implements ContractInterface {
    private final Gson gson=new Gson();
    private final DIDClient didClient;

    public DIDContract(DIDClient didClient) {
        this.didClient = didClient;
    }


    @Transaction(intent = Transaction.TYPE.SUBMIT)
    public String registerDID(Context ctx) throws Exception {
        DID did=didClient.createDID();
        byte[] existing = ctx.getStub().getState(did.DID);
        if (existing != null && existing.length > 0) {
            throw new RuntimeException("DID already exists");
        }
        String json = gson.toJson(did.document);
        ctx.getStub().putState(did.DID, json.getBytes());
        return json;
    }



    //更新密鑰
    @Transaction(intent = Transaction.TYPE.SUBMIT)
    public String updateDID(Context ctx, String did, String recoveryKey,String recoveryPrivateKey) throws NoSuchAlgorithmException, InvalidAlgorithmParameterException, NoSuchAlgorithmException {

        byte[] didData = ctx.getStub().getState(did);//document
        if (didData == null || didData.length == 0) {
            throw new RuntimeException("DID 不存在");
        }
        DIDDocument didDoc = gson.fromJson(new String(didData, StandardCharsets.UTF_8), DIDDocument.class);
        if(!didClient.verifyRecoveryKey(recoveryKey,recoveryPrivateKey)){
            throw new RuntimeException("密鑰為空");
        }
        if(!didDoc.getRecovery().getPublicKey().equals(recoveryKey)){
            throw new RuntimeException("紀錄不匹配");
        }

        //生成
        KeyPair newKeyPair=didClient.generateKeyPair();
        KeyPairs newKeyPairs = new KeyPairs();
        newKeyPairs.setType("RSA");
        newKeyPairs.setPublicKey(didClient.encodePublicKey(newKeyPair));
        newKeyPairs.setPrivateKey(didClient.encodePrivateKey(newKeyPair));

        didDoc.getAuthentication().setPublicKey(newKeyPairs.getPublicKey());
        didDoc.setUpdated(Instant.now().toString());

        String newJson = gson.toJson(didDoc);
        ctx.getStub().putState(did, newJson.getBytes(StandardCharsets.UTF_8));

        return newJson;

    }
    //註銷
    @Transaction(intent = Transaction.TYPE.SUBMIT)
    public void revoke(Context ctx,String did,String recoveryKey,String recoveryPrivateKey) throws NoSuchAlgorithmException, InvalidAlgorithmParameterException, NoSuchAlgorithmException {
        byte[] didData=ctx.getStub().getState(did);//
        if(didData ==null || didData.length==0){
            throw new RuntimeException("DID不存在");
        }

        DIDDocument didDoc = gson.fromJson(new String(didData, StandardCharsets.UTF_8), DIDDocument.class);
        //驗證
        if(!didClient.verifyRecoveryKey(recoveryKey,recoveryPrivateKey)){
            throw new RuntimeException("恢復密鑰不得為空");
        }
        if(!didDoc.getRecovery().getPublicKey().equals(recoveryKey)){
            throw new RuntimeException("密鑰不匹配");
        }

        ctx.getStub().delState(did);
    }




   @Transaction(intent = Transaction.TYPE.EVALUATE)
    public String queryDID(Context ctx, String did) {
    byte[] existing = ctx.getStub().getState(did);
    if (existing == null || existing.length == 0) {
        throw new RuntimeException("DID does not exist");
    }
    return new String(existing);
}








}
