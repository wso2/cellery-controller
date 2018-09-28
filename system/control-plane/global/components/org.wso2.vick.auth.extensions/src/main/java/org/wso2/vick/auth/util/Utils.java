package org.wso2.vick.auth.util;

import org.apache.commons.lang.StringUtils;

public class Utils {

    public static final String OPENID_IDP_ENTITY_ID = "IdPEntityId";

    private Utils() {

    }

    public static boolean isSignedJWT(String jwtToTest) {
        // Signed JWT token contains 3 base64 encoded components separated by periods.
        return StringUtils.countMatches(jwtToTest, ".") == 2;
    }
}
