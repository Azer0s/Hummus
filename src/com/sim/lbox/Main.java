package com.sim.lbox;

import javafx.util.Pair;

import java.io.File;
import java.io.IOException;
import java.net.URISyntaxException;
import java.nio.file.Paths;
import java.util.Scanner;
import javax.script.ScriptEngineManager;
import javax.script.ScriptEngine;
import javax.script.ScriptException;
import java.util.regex.Matcher;

public class Main {

    private ScriptEngineManager mgr = new ScriptEngineManager();
    private ScriptEngine engine = mgr.getEngineByName("JavaScript");

    public static void main(String[] args) {
	// write your code here
        Scanner sc = new Scanner(System.in);

        if (args.length != 0){
            if (new File(args[0]).exists()){
                new FileInterpreter().interpreteAllLine(args[0],new Main());
                try {
                    System.in.read();
                } catch (IOException e) {
                    // ignored
                }
                System.exit(0);
            }else try {
                if(new File(Paths.get(Main.class.getProtectionDomain().getCodeSource().getLocation().toURI().getPath().toString(),args[0]).toString()).exists()){
                    new FileInterpreter().interpreteAllLine(Paths.get(Main.class.getProtectionDomain().getCodeSource().getLocation().toURI().getPath().toString(),args[0]).toString(),new Main());
                    System.in.read();
                    System.exit(0);
                }
            } catch (Exception e) {
                System.out.println(e.getMessage());
            }

            System.out.println("File " + args[0] + "does not exist!");
            try {
                System.in.read();
            } catch (IOException e) {
                // ignored
            }
            System.exit(-1);
        }

        Main m = new Main();
        while(true){
            System.out.print(">");
            String val = sc.nextLine();
            System.out.println(m.interpretLine(val).getKey());
        }
    }

    public Pair<String,Boolean> interpretLine(String line){

        if (line == "exit"){
            System.exit(0);
        }

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
                    calculation = calculation.replace(expectedArgs[i], interpretLine(args[i]).getKey());
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
            String result = interpretLine(m.group(2)).getKey();
            Cache.getInstance().variables.put(m.group(1),result);
            return new Pair<>(m.group(1) + " is " + result,false);
        }else if(Cache.getInstance().variables.containsKey(line)){
            return new Pair<>(Cache.getInstance().variables.get(line),true);
        }else if(line.matches("\\d*") || line.matches("\\w*")){
            return new Pair<>(line,false);
        }

        return new Pair<>("Expression " + line + " is invalid!",true);
    }
}
