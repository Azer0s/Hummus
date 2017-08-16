package com.sim.lbox;

import com.sun.org.apache.xpath.internal.operations.Bool;
import javafx.util.Pair;

import java.util.HashMap;
import java.util.Scanner;
import javax.script.ScriptEngineManager;
import javax.script.ScriptEngine;
import javax.script.ScriptException;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class Main {

    private ScriptEngineManager mgr = new ScriptEngineManager();
    private ScriptEngine engine = mgr.getEngineByName("JavaScript");

    public static void main(String[] args) {
	// write your code here
        Scanner sc = new Scanner(System.in);
        Main m = new Main();
        while(true){
            System.out.print(">");
            String val = sc.nextLine();
            System.out.println(m.InterpretLine(val).getKey());
        }
    }

    public Pair<String,Boolean> InterpretLine(String line){
        if (Cache.getInstance().lExpression.matcher(line).matches()){
            Matcher m = Cache.getInstance().lExpression.matcher(line);
            m.matches();
            String name = m.group(1);
            String[] arguments = m.group(2).split(",");
            String process = m.group(3);

            Cache.getInstance().expressions.put(name,new LExpression(arguments,process));

            return new Pair<>(name + " is " + process + " with arguments: " + m.group(2), false);
        }else if(Cache.getInstance().lCall.matcher(line).matches()){
            Matcher m = Cache.getInstance().lCall.matcher(line);
            m.matches();

            if (Cache.getInstance().expressions.containsKey(m.group(1))){
                String[] args = Cache.getInstance().argSplitter.split(m.group(2));
                String[] expectedArgs = Cache.getInstance().expressions.get(m.group(1)).input;

                if (args.length != expectedArgs.length){
                    return new Pair<>("Invalid amount of arguments!", true);
                }

                String calculation = Cache.getInstance().expressions.get(m.group(1)).calculation;
                for (int i = 0; i < args.length; i++){
                    calculation = calculation.replace(expectedArgs[i],InterpretLine(args[i]).getKey());
                }

                try {
                    return new Pair<String,Boolean>(engine.eval(calculation).toString(),true);
                } catch (ScriptException e) {
                    return new Pair<>(e.getMessage(),true);
                }
            }
        }else if(Cache.getInstance().assignment.matcher(line).matches()){
            Matcher m = Cache.getInstance().assignment.matcher(line);
            m.matches();
            String result = InterpretLine(m.group(2)).getKey();
            Cache.getInstance().variables.put(m.group(1),result);
            return new Pair<>(m.group(1) + " is " + result,false);
        }else if(Cache.getInstance().variables.containsKey(line)){
            return new Pair<>(Cache.getInstance().variables.get(line),false);
        }else if(line.matches("\\d*") || line.matches("\\w*")){
            return new Pair<>(line,false);
        }

        return new Pair<>("Expression " + line + " is invalid!",true);
    }
}
