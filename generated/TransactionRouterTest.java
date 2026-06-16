/* Generated from 'TransactionRouterTest.nrx' 16 Jun 2026 22:47:10 [v5.10] */
/* Options: Annotations Binary Decimal Format Implicituses Java Logo Replace Trace2 Verbose3 */
package com.factory.routing;
import java.sql.SQLException;


public class TransactionRouterTest{
 private static final java.lang.String $0="TransactionRouterTest.nrx";
 
 @SuppressWarnings("unchecked") 
 
 public static void main(java.lang.String args[]){
  com.factory.routing.FuzzInput stringBounds[];
  com.factory.routing.FuzzInput dbPathBounds[];
  com.factory.routing.FuzzInput doubleBounds[];
  com.factory.routing.FuzzInput rexxBounds[];
  java.util.ArrayList recordBounds;
  com.factory.routing.FuzzInput tsVal=null;
  com.factory.routing.FuzzInput nameVal=null;
  com.factory.routing.FuzzInput valVal=null;
  com.factory.routing.TransactionRecord rec=null;
  int recIsCounter=0;
  java.lang.String recExpected=null;
  int counterCount=0;
  com.factory.routing.FuzzInput routeTransaction_val1=null;
  com.factory.routing.FuzzInput routeTransaction_val1_input=null;
  java.lang.String routeTransaction_p1=null;
  java.lang.Object routeTransaction_val2=null;
  com.factory.routing.RecordFuzzInput routeTransaction_val2_input=null;
  com.factory.routing.TransactionRecord routeTransaction_p2=null;
  int routeTransaction_isCounter=0;
  java.util.ArrayList routeTransaction_expectedExs=null;
  java.lang.Throwable routeTransaction_ex=null;
  java.lang.Throwable routeTransaction_caught=null;
  com.factory.routing.FuzzInput initRoutingTable_val1=null;
  com.factory.routing.FuzzInput initRoutingTable_val1_input=null;
  java.lang.String initRoutingTable_p1=null;
  int initRoutingTable_isCounter=0;
  java.util.ArrayList initRoutingTable_expectedExs=null;
  java.lang.Throwable initRoutingTable_ex=null;
  java.lang.Throwable initRoutingTable_caught=null;
  com.factory.routing.FuzzInput getTransactionCount_val1=null;
  com.factory.routing.FuzzInput getTransactionCount_val1_input=null;
  java.lang.String getTransactionCount_p1=null;
  com.factory.routing.FuzzInput getTransactionCount_val2=null;
  com.factory.routing.FuzzInput getTransactionCount_val2_input=null;
  java.lang.String getTransactionCount_p2=null;
  int getTransactionCount_isCounter=0;
  java.util.ArrayList getTransactionCount_expectedExs=null;
  java.lang.Throwable getTransactionCount_ex=null;
  java.lang.Throwable getTransactionCount_caught=null;
  netrexx.lang.RexxIO.Say("=== [Phase III] Starting Boundary Input Exhaustion Test for TransactionRouter ===");
  stringBounds=new com.factory.routing.FuzzInput[]{new com.factory.routing.FuzzInput("normal_string_test",0,(java.lang.String)null),new com.factory.routing.FuzzInput("cpu_usage",0,(java.lang.String)null),new com.factory.routing.FuzzInput("standard_channel",0,(java.lang.String)null),new com.factory.routing.FuzzInput("\' OR \'1\'=\'1",1,"java.sql.SQLException"),new com.factory.routing.FuzzInput("\'; DROP TABLE system_metrics; --",1,"java.sql.SQLException"),new com.factory.routing.FuzzInput("\' UNION SELECT null, null, null --",1,"java.sql.SQLException"),new com.factory.routing.FuzzInput("generated/metrics_test.db",0,(java.lang.String)null),new com.factory.routing.FuzzInput("metrics.db",0,(java.lang.String)null),new com.factory.routing.FuzzInput("routing.db",0,(java.lang.String)null),new com.factory.routing.FuzzInput("/etc/passwd",1,"java.io.IOException"),new com.factory.routing.FuzzInput("../../../etc/passwd",1,"java.io.IOException"),new com.factory.routing.FuzzInput("C:\\Windows\\win.ini",1,"java.io.IOException"),new com.factory.routing.FuzzInput("normal_string",0,(java.lang.String)null),new com.factory.routing.FuzzInput("",1,"java.lang.IllegalArgumentException"),new com.factory.routing.FuzzInput((java.lang.String)null,1,"java.lang.IllegalArgumentException"),new com.factory.routing.FuzzInput("   ",1,"java.lang.IllegalArgumentException")};
  dbPathBounds=new com.factory.routing.FuzzInput[]{new com.factory.routing.FuzzInput("generated/transactionrouter_test.db",0,(java.lang.String)null),new com.factory.routing.FuzzInput(":memory:",0,(java.lang.String)null),new com.factory.routing.FuzzInput((java.lang.String)null,1,"java.lang.IllegalArgumentException")};
  doubleBounds=new com.factory.routing.FuzzInput[]{new com.factory.routing.FuzzInput("0",0,(java.lang.String)null),new com.factory.routing.FuzzInput("1",0,(java.lang.String)null),new com.factory.routing.FuzzInput("-1",0,(java.lang.String)null),new com.factory.routing.FuzzInput("1.5",0,(java.lang.String)null),new com.factory.routing.FuzzInput("-45.2",0,(java.lang.String)null),new com.factory.routing.FuzzInput("100.0",0,(java.lang.String)null),new com.factory.routing.FuzzInput("1.7976931348623157e+308",0,(java.lang.String)null),new com.factory.routing.FuzzInput("-1.7976931348623157e+308",0,(java.lang.String)null),new com.factory.routing.FuzzInput("1.7976931348623159e+308",1,"java.lang.NumberFormatException"),new com.factory.routing.FuzzInput("-1.7976931348623159e+308",1,"java.lang.NumberFormatException")};
  rexxBounds=new com.factory.routing.FuzzInput[]{new com.factory.routing.FuzzInput("normal",0,(java.lang.String)null),new com.factory.routing.FuzzInput("",0,(java.lang.String)null),new com.factory.routing.FuzzInput("0",0,(java.lang.String)null),new com.factory.routing.FuzzInput("1",0,(java.lang.String)null),new com.factory.routing.FuzzInput("-1",0,(java.lang.String)null),new com.factory.routing.FuzzInput("100",0,(java.lang.String)null),new com.factory.routing.FuzzInput("-50",0,(java.lang.String)null),new com.factory.routing.FuzzInput("999",0,(java.lang.String)null),new com.factory.routing.FuzzInput("2147483647",0,(java.lang.String)null),new com.factory.routing.FuzzInput("-2147483648",0,(java.lang.String)null),new com.factory.routing.FuzzInput("2147483648",1,"java.lang.NumberFormatException"),new com.factory.routing.FuzzInput("-2147483649",1,"java.lang.NumberFormatException")};
  recordBounds=new java.util.ArrayList();
  recordBounds.add((java.lang.Object)(new com.factory.routing.RecordFuzzInput((com.factory.routing.TransactionRecord)null,1,"java.lang.IllegalArgumentException")));
  {int $3=0;com.factory.routing.FuzzInput[] $2=new com.factory.routing.FuzzInput[stringBounds.length];synchronized(stringBounds){for(;;){if($3==$2.length)break;$2[$3]=stringBounds[stringBounds.length-1-$3];$3++;}}tsVal:for(;;){if(--$3<0)break;tsVal=(com.factory.routing.FuzzInput)$2[$3];
   {int $6=0;com.factory.routing.FuzzInput[] $5=new com.factory.routing.FuzzInput[stringBounds.length];synchronized(stringBounds){for(;;){if($6==$5.length)break;$5[$6]=stringBounds[stringBounds.length-1-$6];$6++;}}nameVal:for(;;){if(--$6<0)break;nameVal=(com.factory.routing.FuzzInput)$5[$6];
    {int $9=0;com.factory.routing.FuzzInput[] $8=new com.factory.routing.FuzzInput[doubleBounds.length];synchronized(doubleBounds){for(;;){if($9==$8.length)break;$8[$9]=doubleBounds[doubleBounds.length-1-$9];$9++;}}valVal:for(;;){if(--$9<0)break;valVal=(com.factory.routing.FuzzInput)$8[$9];
     rec=new com.factory.routing.TransactionRecord();
     if (tsVal.val!=null) 
      rec.priority=tsVal.val;
     if (nameVal.val!=null) 
      rec.receiver=nameVal.val;
     if (tsVal.val!=null) 
      rec.sender=tsVal.val;
     if (valVal.val!=null) 
      rec.amount=new netrexx.lang.Rexx(valVal.val);
     if (nameVal.val!=null) 
      rec.txId=nameVal.val;
     recIsCounter=(tsVal.isCounter|nameVal.isCounter)|valVal.isCounter;
     recExpected=(java.lang.String)null;
     if (tsVal.isCounter!=0) 
      {
       recExpected=tsVal.expected;
      }
     else 
      if (nameVal.isCounter!=0) 
       {
        recExpected=nameVal.expected;
       }
      else 
       if (valVal.isCounter!=0) 
        {
         recExpected=valVal.expected;
        }
     counterCount=(tsVal.isCounter+nameVal.isCounter)+valVal.isCounter;
     if ((counterCount==0)|(counterCount==1)) 
      {
       recordBounds.add((java.lang.Object)(new com.factory.routing.RecordFuzzInput(rec,recIsCounter,recExpected)));
      }
     }
    }/*valVal*/
    }
   }/*nameVal*/
   }
  }/*tsVal*/
  if (((((stringBounds!=null)&(dbPathBounds!=null))&(doubleBounds!=null))&(rexxBounds!=null))&(recordBounds!=null)) 
   netrexx.lang.RexxIO.Say("ok");
  netrexx.lang.RexxIO.Say("Testing method routeTransaction...");
  {int $12=0;com.factory.routing.FuzzInput[] $11=new com.factory.routing.FuzzInput[dbPathBounds.length];synchronized(dbPathBounds){for(;;){if($12==$11.length)break;$11[$12]=dbPathBounds[dbPathBounds.length-1-$12];$12++;}}routeTransaction_val1:for(;;){if(--$12<0)break;routeTransaction_val1=(com.factory.routing.FuzzInput)$11[$12];
   routeTransaction_val1_input=routeTransaction_val1;
   routeTransaction_p1=(java.lang.String)null;
   {int $15=0;java.lang.Object[] $14=new java.lang.Object[recordBounds.size()];synchronized(recordBounds){java.util.Iterator $13=recordBounds.iterator();for(;;){if($15==$14.length)break;$14[$15]=$13.next();$15++;}}routeTransaction_val2:for(;;){if(--$15<0)break;routeTransaction_val2=(java.lang.Object)$14[$15];
    routeTransaction_val2_input=(com.factory.routing.RecordFuzzInput)routeTransaction_val2;
    routeTransaction_p2=(com.factory.routing.TransactionRecord)null;
    routeTransaction_isCounter=0;
    routeTransaction_expectedExs=new java.util.ArrayList();
    if (routeTransaction_val1_input.isCounter!=0) 
     {
      routeTransaction_isCounter=1;
      routeTransaction_expectedExs.add((java.lang.Object)routeTransaction_val1_input.expected);
     }
    if (routeTransaction_val2_input.isCounter!=0) 
     {
      routeTransaction_isCounter=1;
      routeTransaction_expectedExs.add((java.lang.Object)routeTransaction_val2_input.expected);
     }
    routeTransaction_ex=(java.lang.Throwable)null;
    {try{
     if (1==0) 
      com.factory.routing.TransactionRouterTest.dummySignal();
     if (routeTransaction_val1_input.val!=null) 
      routeTransaction_p1=routeTransaction_val1_input.val;
     routeTransaction_p2=routeTransaction_val2_input.rec;
     com.factory.routing.TransactionRouter.routeTransaction(routeTransaction_p1,routeTransaction_p2);
    }
    catch (java.lang.Throwable $16){routeTransaction_caught=$16;
     routeTransaction_ex=routeTransaction_caught;
    }}
    com.factory.routing.TransactionRouterTest.assertResult("routeTransaction",routeTransaction_isCounter,routeTransaction_expectedExs,routeTransaction_ex);
    }
   }/*routeTransaction_val2*/
   }
  }/*routeTransaction_val1*/
  netrexx.lang.RexxIO.Say("  Method routeTransaction boundary exhaustion completed.");
  netrexx.lang.RexxIO.Say("Testing method initRoutingTable...");
  {int $19=0;com.factory.routing.FuzzInput[] $18=new com.factory.routing.FuzzInput[dbPathBounds.length];synchronized(dbPathBounds){for(;;){if($19==$18.length)break;$18[$19]=dbPathBounds[dbPathBounds.length-1-$19];$19++;}}initRoutingTable_val1:for(;;){if(--$19<0)break;initRoutingTable_val1=(com.factory.routing.FuzzInput)$18[$19];
   initRoutingTable_val1_input=initRoutingTable_val1;
   initRoutingTable_p1=(java.lang.String)null;
   initRoutingTable_isCounter=0;
   initRoutingTable_expectedExs=new java.util.ArrayList();
   if (initRoutingTable_val1_input.isCounter!=0) 
    {
     initRoutingTable_isCounter=1;
     initRoutingTable_expectedExs.add((java.lang.Object)initRoutingTable_val1_input.expected);
    }
   initRoutingTable_ex=(java.lang.Throwable)null;
   {try{
    if (1==0) 
     com.factory.routing.TransactionRouterTest.dummySignal();
    if (initRoutingTable_val1_input.val!=null) 
     initRoutingTable_p1=initRoutingTable_val1_input.val;
    com.factory.routing.TransactionRouter.initRoutingTable(initRoutingTable_p1);
   }
   catch (java.lang.Throwable $20){initRoutingTable_caught=$20;
    initRoutingTable_ex=initRoutingTable_caught;
   }}
   com.factory.routing.TransactionRouterTest.assertResult("initRoutingTable",initRoutingTable_isCounter,initRoutingTable_expectedExs,initRoutingTable_ex);
   }
  }/*initRoutingTable_val1*/
  netrexx.lang.RexxIO.Say("  Method initRoutingTable boundary exhaustion completed.");
  netrexx.lang.RexxIO.Say("Testing method getTransactionCount...");
  {int $23=0;com.factory.routing.FuzzInput[] $22=new com.factory.routing.FuzzInput[dbPathBounds.length];synchronized(dbPathBounds){for(;;){if($23==$22.length)break;$22[$23]=dbPathBounds[dbPathBounds.length-1-$23];$23++;}}getTransactionCount_val1:for(;;){if(--$23<0)break;getTransactionCount_val1=(com.factory.routing.FuzzInput)$22[$23];
   getTransactionCount_val1_input=getTransactionCount_val1;
   getTransactionCount_p1=(java.lang.String)null;
   {int $26=0;com.factory.routing.FuzzInput[] $25=new com.factory.routing.FuzzInput[stringBounds.length];synchronized(stringBounds){for(;;){if($26==$25.length)break;$25[$26]=stringBounds[stringBounds.length-1-$26];$26++;}}getTransactionCount_val2:for(;;){if(--$26<0)break;getTransactionCount_val2=(com.factory.routing.FuzzInput)$25[$26];
    getTransactionCount_val2_input=getTransactionCount_val2;
    getTransactionCount_p2=(java.lang.String)null;
    getTransactionCount_isCounter=0;
    getTransactionCount_expectedExs=new java.util.ArrayList();
    if (getTransactionCount_val1_input.isCounter!=0) 
     {
      getTransactionCount_isCounter=1;
      getTransactionCount_expectedExs.add((java.lang.Object)getTransactionCount_val1_input.expected);
     }
    if (getTransactionCount_val2_input.isCounter!=0) 
     {
      getTransactionCount_isCounter=1;
      getTransactionCount_expectedExs.add((java.lang.Object)getTransactionCount_val2_input.expected);
     }
    getTransactionCount_ex=(java.lang.Throwable)null;
    {try{
     if (1==0) 
      com.factory.routing.TransactionRouterTest.dummySignal();
     if (getTransactionCount_val1_input.val!=null) 
      getTransactionCount_p1=getTransactionCount_val1_input.val;
     if (getTransactionCount_val2_input.val!=null) 
      getTransactionCount_p2=getTransactionCount_val2_input.val;
     com.factory.routing.TransactionRouter.getTransactionCount(getTransactionCount_p1,getTransactionCount_p2);
    }
    catch (java.lang.Throwable $27){getTransactionCount_caught=$27;
     getTransactionCount_ex=getTransactionCount_caught;
    }}
    com.factory.routing.TransactionRouterTest.assertResult("getTransactionCount",getTransactionCount_isCounter,getTransactionCount_expectedExs,getTransactionCount_ex);
    }
   }/*getTransactionCount_val2*/
   }
  }/*getTransactionCount_val1*/
  netrexx.lang.RexxIO.Say("  Method getTransactionCount boundary exhaustion completed.");
  netrexx.lang.RexxIO.Say("=== [Phase III] Boundary Input Exhaustion Test Completed successfully! ===");
  return;}
 
 
 @SuppressWarnings("unchecked") 
 
 public static void assertResult(java.lang.String methodName,int isCounter,java.util.ArrayList expectedExs,java.lang.Throwable ex){
  java.lang.String thrownEx=null;
  int npeExpected=0;
  int i=0;
  int matched=0;
  java.lang.String expEx=null;
  java.lang.Class expectedClass=null;
  int idx=0;
  java.lang.String exp=null;
  if (isCounter!=0) 
   {
    if (ex==null) 
     {
      netrexx.lang.RexxIO.Say("Assertion failure in "+methodName+": counter-example bypassed validation (no exception thrown)");
      java.lang.System.exit(1);
     }
    thrownEx=ex.getClass().getName();
    if (thrownEx.equals("java.lang.NullPointerException")) 
     {
      npeExpected=0;
      {int $28=(expectedExs.size())-1;i=0;i:for(;i<=$28;i++){
       if ((((java.lang.String)(expectedExs.get(i)))).equals("java.lang.NullPointerException")) 
        npeExpected=1;
       }
      }/*i*/
      if (npeExpected==0) 
       {
        netrexx.lang.RexxIO.Say("Assertion failure in "+methodName+": ungraceful crash (NullPointerException)");
        ex.printStackTrace();
        java.lang.System.exit(1);
       }
     }
    matched=0;
    {int $29=(expectedExs.size())-1;i=0;i:for(;i<=$29;i++){
     expEx=(java.lang.String)(expectedExs.get(i));
     {try{
      expectedClass=java.lang.Class.forName(expEx);
      if (expectedClass.isInstance((java.lang.Object)ex)) 
       matched=1;
     }
     catch (java.lang.ClassNotFoundException $30){
      ;
     }}
     }
    }/*i*/
    if (matched==0) 
     {
      netrexx.lang.RexxIO.Say("Assertion failure in "+methodName+": caught "+thrownEx+" ("+ex.getMessage()+") but none of the expected exceptions matched.");
      netrexx.lang.RexxIO.Say("Expected exceptions for "+methodName+":");
      {int $31=(expectedExs.size())-1;idx=0;idx:for(;idx<=$31;idx++){
       exp=(java.lang.String)(expectedExs.get(idx));
       if (exp==null) 
        exp="null";
       netrexx.lang.RexxIO.Say("  - "+exp);
       }
      }/*idx*/
      java.lang.System.exit(1);
     }
   }
  else 
   {
    if (ex!=null) 
     {
      netrexx.lang.RexxIO.Say("Assertion failure in "+methodName+": happy path regression (unexpected exception "+ex.getClass().getName()+": "+ex.getMessage()+")");
      ex.printStackTrace();
      java.lang.System.exit(1);
     }
   }
  return;}
 
 
 @SuppressWarnings("unchecked") 
 
 private static void dummySignal() throws java.lang.Throwable{
  ;
  return;}
 
 
 private TransactionRouterTest(){return;}
 }


class FuzzInput{
 private static final java.lang.String $0="TransactionRouterTest.nrx";
 /* properties public */
 public java.lang.String val;
 public int isCounter;
 public java.lang.String expected;
 
 @SuppressWarnings("unchecked") 
 
 public FuzzInput(java.lang.String aVal,int aIsCounter,java.lang.String aExpected){super();
  this.val=aVal;
  this.isCounter=aIsCounter;
  this.expected=aExpected;
  return;}
 
 }


class RecordFuzzInput{
 private static final java.lang.String $0="TransactionRouterTest.nrx";
 /* properties public */
 public com.factory.routing.TransactionRecord rec;
 public int isCounter;
 public java.lang.String expected;
 
 @SuppressWarnings("unchecked") 
 
 public RecordFuzzInput(com.factory.routing.TransactionRecord aRec,int aIsCounter,java.lang.String aExpected){super();
  this.rec=aRec;
  this.isCounter=aIsCounter;
  this.expected=aExpected;
  return;}
 
 }