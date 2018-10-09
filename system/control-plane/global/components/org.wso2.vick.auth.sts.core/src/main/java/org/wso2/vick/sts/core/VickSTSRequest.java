package org.wso2.vick.sts.core;

import java.util.List;

public class VickSTSRequest {

    /**
     * Identifier of the workload initiating the STS request.
     */
    private String source;

    private List<String> scopes;

    private List<String> audiences;

    public String getSource() {

        return source;
    }

    public void setSource(String source) {

        this.source = source;
    }

    public List<String> getScopes() {

        return scopes;
    }

    public void setScopes(List<String> scopes) {

        this.scopes = scopes;
    }

    public List<String> getAudiences() {

        return audiences;
    }

    public void setAudiences(List<String> audiences) {

        this.audiences = audiences;
    }
}
