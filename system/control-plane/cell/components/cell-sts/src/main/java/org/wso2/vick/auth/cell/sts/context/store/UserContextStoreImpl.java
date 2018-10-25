/*
 *  Copyright (c) 2018 WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
 *
 *  WSO2 Inc. licenses this file to you under the Apache License,
 *  Version 2.0 (the "License"); you may not use this file except
 *  in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing,
 *  software distributed under the License is distributed on an
 *  "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 *  KIND, either express or implied.  See the License for the
 *  specific language governing permissions and limitations
 *  under the License.
 */
package org.wso2.vick.auth.cell.sts.context.store;

import com.google.common.cache.Cache;
import com.google.common.cache.CacheBuilder;
import org.apache.commons.lang.StringUtils;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.Map;
import java.util.concurrent.TimeUnit;

public class UserContextStoreImpl implements UserContextStore {

    private static final String USER_CONTEXT_EXPIRY_IN_SECONDS = "USER_CONTEXT_EXPIRY_SECONDS";
    private static final long DEFAULT_EXPIRY_IN_SECONDS = 300L;
    private static final Logger log = LoggerFactory.getLogger(UserContextStoreImpl.class);

    private Map<String, String> userContextMap;

    public UserContextStoreImpl() {

        log.info("User Context expiry set to {} seconds." , getUserContextExpiry());
        Cache<String, String> cache = CacheBuilder.newBuilder()
                .expireAfterWrite(getUserContextExpiry(), TimeUnit.SECONDS)
                .removalListener(removalNotification -> {
                    log.debug("Stored user context was removed: " + removalNotification);
                })
                .build();

        userContextMap = cache.asMap();
    }

    private long getUserContextExpiry() {

        long expiryTimeInSecs = DEFAULT_EXPIRY_IN_SECONDS;
        String expiryConfigValue = System.getenv(USER_CONTEXT_EXPIRY_IN_SECONDS);
        if (StringUtils.isNotBlank(expiryConfigValue)) {
            try {
                expiryTimeInSecs = Long.parseLong(expiryConfigValue);
            } catch (NumberFormatException ex) {
                log.warn("Invalid value '{}' provided for user context store expiry. Using default value: {}",
                        expiryConfigValue, expiryTimeInSecs);
            }
        }
        return expiryTimeInSecs;
    }

    @Override
    public String get(String contextId) {
        return userContextMap.get(contextId);
    }

    @Override
    public void put(String contextId, String context) {
        userContextMap.put(contextId, context);
    }

    @Override
    public boolean containsKey(String contextId) {
        return userContextMap.containsKey(contextId);
    }


}
