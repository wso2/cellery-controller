package org.wso2.vick.auth.km;

import org.wso2.carbon.apimgt.keymgt.APIKeyMgtException;
import org.wso2.carbon.apimgt.keymgt.handlers.DefaultKeyValidationHandler;
import org.wso2.carbon.apimgt.keymgt.service.TokenValidationContext;
import org.wso2.vick.auth.util.Utils;

public class VickExtendedKeyValidationHandler extends DefaultKeyValidationHandler {

    @Override
    public boolean validateSubscription(TokenValidationContext validationContext) throws APIKeyMgtException {

        if (Utils.isSignedJWT(validationContext.getAccessToken())) {
            // We are skipping subscription validation for the moment.
            return true;
        }

        return super.validateSubscription(validationContext);
    }

}
