package com.sim.lbox;

import java.util.HashMap;
import java.util.regex.Pattern;

import static java.util.regex.Pattern.*;

/**
 * Created by ariel on 16.08.2017.
 */
class Cache {

    Pattern argSplitter = compile(",(?=([^\\(]*\\([^\\\"]*\\))*[^\\)]*$)");
    Pattern lExpression = compile("(\\w*) *:= *\\(([\\w,]*)\\)\\.\\((.*)\\)");
    Pattern assignment = compile("(\\w*) *:= *(.*)");
    Pattern anonymous = compile("\\(([\\w,]*)\\)\\.\\((.*)\\)\\.\\((.*)\\)");
    Pattern ifCondition = compile("(.+)\\?(.+):(.+)");
    Pattern lCall = compile("(\\w*)\\((.*)\\)");
    boolean rec = false;
    HashMap<String,LExpression> expressions = new HashMap<>();
    HashMap<String,String> variables = new HashMap<>();
    private static Cache instance;

    private Cache () {}
    static Cache getInstance() {
        if (Cache.instance == null) {
            Cache.instance = new Cache ();
        }
        return Cache.instance;
    }
}
