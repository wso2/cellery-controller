/*
 * Copyright (c) 2018, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
 *
 * WSO2 Inc. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 *
 */
package org.wso2.vick.auth.cell.sts;

import com.google.gson.Gson;
import com.google.gson.GsonBuilder;
import org.apache.commons.lang.StringUtils;
import org.json.simple.JSONObject;
import org.wso2.vick.auth.cell.sts.service.VickCellSTSException;

import java.util.Map;

public class CellStsUtils {

    private static final String CELL_NAME_ENV_VARIABLE = "CELL_NAME";

    public static String getMyCellName() throws VickCellSTSException {
        // For now we pick the cell name from the environment variable.
        String cellName = System.getenv(CELL_NAME_ENV_VARIABLE);
        if (StringUtils.isBlank(cellName)) {
            throw new VickCellSTSException("Environment variable '" + CELL_NAME_ENV_VARIABLE + "' is empty.");
        }
        return cellName;
    }

    public static boolean isWorkloadExternalToVick(String destinationWorkloadName) {
        // For now the only way to check whether the destination is outside of VICK/Cell mesh is by checking whether the
        // destination workload name does not comply to the format <cell-name>--<service_name>
        // Eg: hr--employee-service
        // Once we find a smarter way to check whether the target is outside of VICK we can replace this not-so-smart
        // logic.
        return !StringUtils.contains(destinationWorkloadName, "--");
    }

    public static String getPrettyPrintJson(Map<String, String> attributes) {
        JSONObject configJson = new JSONObject();
        attributes.forEach((key, value) -> configJson.put(key, value));

        Gson gson = new GsonBuilder().setPrettyPrinting().create();
        return gson.toJson(configJson);
    }

    /**
     * Returns the issuer name.
     * @param cellName  Name of the cell.
     * @return Issuer name of the respective cell.
     */
    public static String getIssuerName(String cellName) {

        return cellName + "--sts-service";
    }

}
