package DIDAPI;


import com.ibm.cloud.sdk.core.http.HttpStatus;
import org.hyperledger.fabric.gateway.Contract;
import org.hyperledger.fabric.gateway.ContractException;
import org.hyperledger.fabric.gateway.Network;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import javax.annotation.Resource;
import java.util.concurrent.TimeoutException;

@RestController
@RequestMapping("/did")
public class DIDController {

    @Resource
    private Contract contract;

    @Resource
    private Network network;

    @PostMapping("/registerDID")
    public ResponseEntity<String> registerDID() throws ContractException, InterruptedException, TimeoutException {
        byte[] result =contract.submitTransaction("registerDID");
        return ResponseEntity.ok(new String(result));
    }
    @PostMapping("/updateDID")
    public ResponseEntity<String> updateDID(@RequestParam String did,@RequestParam String recoveryKey,@RequestParam String recoveryPrivateKey) throws ContractException, InterruptedException, TimeoutException {
        byte[] result=contract.submitTransaction("updateDID",did,recoveryKey,recoveryPrivateKey);
        return ResponseEntity.ok(new String(result));
    }
    @GetMapping("/query/{did}")
    public ResponseEntity<String> queryDID(@PathVariable String did) throws ContractException {
        byte[] result=contract.evaluateTransaction("queryDID",did);
        return ResponseEntity.ok(new String(result));
    }
    @PostMapping("/deleteDID")
    public ResponseEntity<String> deleteDID(@RequestParam String did,@RequestParam String recoveryKey, @RequestParam String recoveryPrivateKey) throws ContractException, InterruptedException, TimeoutException {
        try {
            byte[] result = contract.submitTransaction("revokeDID", did, recoveryKey, recoveryPrivateKey);
            return ResponseEntity.ok("註銷成功");
        }catch(Exception e){
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).body("註銷失敗");
        }
    }

}
