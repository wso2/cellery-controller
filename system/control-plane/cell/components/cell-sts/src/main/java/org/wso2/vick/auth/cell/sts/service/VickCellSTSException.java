package org.wso2.vick.auth.cell.sts.service;

public class VickCellSTSException extends Exception {

    public VickCellSTSException(String message) {

        super(message);
    }

    public VickCellSTSException(Throwable e) {

        super(e);
    }

    public VickCellSTSException(String errorMessage, Throwable e) {

        super(errorMessage, e);
    }
}
