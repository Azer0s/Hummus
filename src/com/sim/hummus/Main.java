package com.sim.hummus;

import javax.script.ScriptEngine;
import javax.script.ScriptEngineManager;
import javax.script.ScriptException;
import java.io.File;
import java.io.FileReader;
import java.io.IOException;
import java.nio.file.Paths;
import java.util.*;
import java.util.regex.Matcher;

import static java.lang.System.exit;
import static java.lang.System.getProperty;
import static java.lang.System.out;

@SuppressWarnings({"ResultOfMethodCallIgnored", "InfiniteLoopStatement"})
public class Main {

    private static ScriptEngineManager mgr = new ScriptEngineManager();
    private static ScriptEngine engine = mgr.getEngineByName("JavaScript");

    public static void main(String[] args) {
        Scanner sc = new Scanner(System.in);

        //Add standard functions
        try {
            ClassLoader classLoader = Main.class.getClassLoader();
            File file = new File(classLoader.getResource("stdlib.js").getFile());
            engine.eval(new FileReader(file));
        } catch (Exception e) {
            e.printStackTrace();
        }

        if (args.length != 0){
            if (new File(args[0]).exists()){
                new FileInterpreter().interpreteAllLine(args[0],new Main());
                try {
                    System.in.read();
                } catch (IOException e) {
                    // ignored
                }
                exit(0);
            }else try {
                if(new File(Paths.get(Main.class.getProtectionDomain().getCodeSource().getLocation().toURI().getPath().toString(),args[0]).toString()).exists()){
                    new FileInterpreter().interpreteAllLine(Paths.get(Main.class.getProtectionDomain().getCodeSource().getLocation().toURI().getPath().toString(),args[0]).toString(),new Main());
                    System.in.read();
                    exit(0);
                }
            } catch (Exception e) {
                out.println(e.getMessage());
            }

            out.println("File " + args[0] + "does not exist!");
            try {
                System.in.read();
            } catch (IOException e) {
                // ignored
            }
            exit(-1);
        }

        Main m = new Main();
        while(true){
            out.print(">");
            String val = sc.nextLine();

            if (val.startsWith("help")){

                if(val.equals("help")){
                    out.println( "- functions\n" +
                                        "- variables\n" +
                                        "- misc\n");
                    continue;
                }

                String help = val.split(" ")[1];

                switch (help) {
                    case "functions":
                        out.println(
                                "Function assignment\n" +
                                        "-------------------\n" +
                                        "name:=(arguments - comma separated).(process)\n\n" +
                                        "Use function\n" +
                                        "-------------------\n" +
                                        "name(arguments - comma seperated)\n\n" +
                                        "Anonymous functions\n" +
                                        "-------------------\n" +
                                        "(arguments - comma separated).(process).(values - comma separated)\n\n" +
                                        "Examples\n" +
                                        "-------------------\n" +
                                        "y:=(x).(x*x)\n" +
                                        "z:=y(y(x))\n" +
                                        "b:=a(true,false)\n" +
                                        "(y,x).(x-y).(1,2)\n" +
                                        "k:=(i).(i ? y(10):y(12))\n"+
                                        "k(true)\n");
                        continue;
                    case "variables":
                        out.println(
                                "Variable assignment\n" +
                                        "-------------------\n" +
                                        "name:=value\n\n" +
                                        "Use variable\n" +
                                        "-------------------\n" +
                                        "function(name)\n\n" +
                                        "Print variable\n" +
                                        "-------------------\n" +
                                        "name\n\n" +
                                        "Examples\n" +
                                        "-------------------\n" +
                                        "x:=2\n" +
                                        "z:=y(x)\n" +
                                        "z\n");
                        continue;
                    case "misc":
                        out.println(
                                "Quit the application\n" +
                                        "-------------------\n" +
                                        "exit\n\n" +
                                        "Clear the console\n" +
                                        "-------------------\n" +
                                        "clear\n");
                        continue;
                    default:
                        out.println("Invalid input!");
                        continue;
                }
            }

            try{
                Pair<String,Boolean> result = m.interpretLine(val);

                if (result != null){
                    out.println(result.getKey());
                }
            }catch (Exception e){
                out.println("Operation " + val + " is invalid!");
            }
        }
    }

    Pair<String,Boolean> interpretLine(String line){

        if (line.equals("exit")){
            exit(0);
        }

        if (line.equals("clear")){
            clearConsole();
            return null;
        }

        if(line.equals("rec")){
            Cache.getInstance().rec = !Cache.getInstance().rec;

            if (Cache.getInstance().rec){
                return new Pair<> ("Recursion is enabled!",false);
            }else{
                return new Pair<> ("Recursion is disabled!",false);
            }
        }

        if (Cache.getInstance().anonymous.matcher(line).matches()){
            Matcher m = Cache.getInstance().anonymous.matcher(line);
            m.matches();
            String calculation = null;
            try {
                calculation = GetCalculation(m.group(2), Cache.getInstance().argSplitter.split(m.group(1)),Cache.getInstance().argSplitter.split(m.group(3)));
            } catch (Exception e) {
                return new Pair<>(e.getMessage(),true);
            }
            return FunctionCalc(calculation);
        }else if (Cache.getInstance().lExpression.matcher(line).matches()){
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
                String calculation;
                try {
                    calculation = GetCalculation(Cache.getInstance().expressions.get(m.group(1)).calculation, Cache.getInstance().expressions.get(m.group(1)).input,Cache.getInstance().argSplitter.split(m.group(2)));
                } catch (Exception e) {
                    return new Pair<>(e.getMessage(),true);
                }

                return FunctionCalc(calculation);
            }
        }else if(Cache.getInstance().assignment.matcher(line).matches()){
            Matcher m = Cache.getInstance().assignment.matcher(line);
            m.matches();
            String result = interpretLine(m.group(2)).getKey();
            Cache.getInstance().variables.put(m.group(1),result);
            return new Pair<>(m.group(1) + " is " + result,false);
        }else if(Cache.getInstance().ifCondition.matcher(line).matches()){
            Matcher m = Cache.getInstance().ifCondition.matcher(line);
            m.matches();
            String condition = interpretLine(m.group(1)).getKey();
            String action1 = m.group(2).trim();
            String action2 = m.group(3).trim();

            if (Boolean.parseBoolean(condition)){
                return interpretLine(action1);
            }else {
                return interpretLine(action2);
            }
        }else if(Cache.getInstance().variables.containsKey(line)){
            return new Pair<>(Cache.getInstance().variables.get(line),true);
        }else if(line.matches("\\d*") || line.matches("\\w*")){
            return new Pair<>(line,false);
        }

        try{
            Object result = engine.eval(line);

            if (result != null){
                return new Pair<>(result.toString(),true);
            }else {
                return new Pair<>("",true);
            }
        }catch (Exception e){
            return new Pair<>("Expression " + line + " is invalid!",true);
        }
    }

    private String GetCalculation(String calculation, String[] expectedArgs,String[] args) throws Exception {
        if (args.length != expectedArgs.length){
            throw new Exception("Invalid amount of arguments!");
        }

        HashMap<String,String> map = new HashMap<>();
        for (int i = 0; i < args.length; i++){
            map.put(expectedArgs[i],interpretLine(args[i]).getKey());
        }

        List<String> keys = new ArrayList<>(map.keySet());
        keys.sort(aStringComparator);

        for (String s: keys) {
            calculation = calculation.replace(s,map.get(s));
        }

        return calculation;
    }

    private Pair<String,Boolean> FunctionCalc(String calculation) {

        if (Cache.getInstance().ifCondition.matcher(calculation).matches()){
            return interpretLine(calculation);
        }

        try {
            return new Pair<>(engine.eval(calculation).toString(),true);
        } catch (ScriptException e) {
            //Don´t do stupid things just because you can...
            if (Cache.getInstance().rec){
                return new Pair<>(interpretLine(calculation).getKey(),true);
            }else {
                if(!calculation.contains("(") && !calculation.contains(")")){
                    return new Pair<>(calculation,true);
                }

                return new Pair<>("Operation " + calculation + " is invalid! It might be recursive!",true);
            }
        }
    }

    private static void clearConsole()
    {
        try
        {
            final String os = getProperty("os.name");

            if (os.contains("Windows"))
            {
                new ProcessBuilder("cmd", "/c", "cls").inheritIO().start().waitFor();
            }
            else
            {
                out.print("\033[H\033[2J");
                out.flush();
            }
        }
        catch (final Exception e)
        {
            //  Handle any exceptions.
        }
    }

    private static final Comparator<String> aStringComparator = (o1, o2) -> {
        //assumed input are strings in the form axxxx
        return Integer.compare(o2.length(), o1.length());
    };

}
