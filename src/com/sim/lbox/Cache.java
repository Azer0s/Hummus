package com.sim.lbox;

import java.util.HashMap;
import java.util.regex.Pattern;

/**
 * Created by ariel on 16.08.2017.
 */
public class Cache {

    public Pattern argSplitter = Pattern.compile(",(?=([^\\(]*\\([^\\\"]*\\))*[^\\)]*$)");
    public Pattern lExpression = Pattern.compile("(\\w*) *:= *\\(([\\w,]*)\\)\\.\\((.*)\\)");
    public Pattern assignment = Pattern.compile("(\\w*) *:= *(.*)");
    public Pattern anonymous = Pattern.compile("\\(([\\w,]*)\\)\\.\\((.*)\\)\\.\\((.*)\\)");
    public Pattern lCall = Pattern.compile("(\\w*)\\((.*)\\)");
    public boolean rec = false;
    public HashMap<String,LExpression> expressions = new HashMap<String,LExpression>();
    public HashMap<String,String> variables = new HashMap<String,String >();
    private static Cache instance;

    private Cache () {}
    public static Cache getInstance () {
        if (Cache.instance == null) {
            Cache.instance = new Cache ();
        }
        return Cache.instance;
    }
}
