package com.DID;

import com.google.common.cache.Cache;
import com.google.common.cache.CacheBuilder;

import java.util.concurrent.TimeUnit;

public class DIDCache {
    private static final Cache<String, String> cache = CacheBuilder.newBuilder()
            .expireAfterWrite(10, TimeUnit.MINUTES)  // 10 分钟过期
            .maximumSize(1000)
            .build();

    public static void put(String did, String json) {
        cache.put(did, json);
    }

    public static boolean contains(String did) {
        return cache.getIfPresent(did) != null;
    }

    public static String get(String did) {
        return cache.getIfPresent(did);
    }
}
