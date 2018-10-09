package org.wso2.vick.auth.exception;

public class VickAuthException extends Exception {

    public VickAuthException(String msg, Throwable throwable) {

        super(msg, throwable);
    }

    public VickAuthException(Throwable throwable) {

        super(throwable);
    }
}
