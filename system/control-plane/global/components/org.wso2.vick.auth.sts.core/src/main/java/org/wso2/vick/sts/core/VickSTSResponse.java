package org.wso2.vick.sts.core;

import org.json.simple.JSONObject;

public class VickSTSResponse {

    private String stsToken;

    public String getStsToken() {

        return stsToken;
    }

    public void setStsToken(String stsToken) {

        this.stsToken = stsToken;
    }

    public String toJson() {

        JSONObject tokenResponse = new JSONObject();
        tokenResponse.put(VickSTSConstants.VickSTSResponse.STS_TOKEN, stsToken);
        return tokenResponse.toString();
    }
}
