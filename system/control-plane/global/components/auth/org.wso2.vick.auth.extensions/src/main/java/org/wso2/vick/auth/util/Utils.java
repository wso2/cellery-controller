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

package org.wso2.vick.auth.util;

import org.apache.commons.lang.StringUtils;

/**
 * This is the Utils class for token validations.
 *
 */
public class Utils {

    public static final String OPENID_IDP_ENTITY_ID = "IdPEntityId";

    private Utils() {

    }

    public static boolean isSignedJWT(String jwtToTest) {
        // Signed JWT token contains 3 base64 encoded components separated by periods.
        return StringUtils.countMatches(jwtToTest, ".") == 2;
    }
}
