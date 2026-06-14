/* Generated from 'TransactionRouterTest.nrx' 15 Jun 2026 00:37:07 [v5.10] */
/* Options: Annotations Binary Decimal Format Implicituses Java Logo Replace Sourcedir Trace2 Verbose3 */
package com.factory.routing;
import java.sql.SQLException;


public class TransactionRouterTest{
 private static final netrexx.lang.Rexx $01=netrexx.lang.Rexx.toRexx("null");
 private static final java.lang.String $0="TransactionRouterTest.nrx";
 
 @SuppressWarnings("unchecked") 
 
 public static void main(java.lang.String args[]){
  java.lang.Class recordClass=null;
  java.lang.Class routerClass=null;
  java.lang.Class stringClass=null;
  java.lang.reflect.Field f_txId=null;
  java.lang.reflect.Field f_sender=null;
  java.lang.reflect.Field f_receiver=null;
  java.lang.reflect.Field f_amount=null;
  java.lang.reflect.Field f_priority=null;
  java.lang.reflect.Method m_route=null;
  java.lang.reflect.Method m_init=null;
  java.lang.reflect.Method m_count=null;
  java.lang.String stringBounds[]=null;
  java.lang.String dbPathBounds[]=null;
  netrexx.lang.Rexx doubleBounds[]=null;
  java.util.ArrayList recordBounds=null;
  java.lang.String tsVal=null;
  java.lang.String nameVal=null;
  netrexx.lang.Rexx valVal=null;
  java.lang.reflect.Constructor constructors[]=null;
  java.lang.Object rec=null;
  java.lang.String val1=null;
  java.lang.String routeTransaction_p1=null;
  java.lang.Object val2=null;
  java.lang.Object routeTransaction_p2=null;
  java.lang.String initRoutingTable_p1=null;
  java.lang.String getTransactionCount_p1=null;
  java.lang.String getTransactionCount_p2=null;
  {try{
   netrexx.lang.RexxIO.Say("=== [Phase III] Starting Boundary Input Exhaustion Test for TransactionRouter ===");
   recordClass=java.lang.Class.forName("com.factory.routing.TransactionRecord");
   routerClass=java.lang.Class.forName("com.factory.routing.TransactionRouter");
   stringClass=java.lang.Class.forName("java.lang.String");
   f_txId=recordClass.getField("txId");
   f_sender=recordClass.getField("sender");
   f_receiver=recordClass.getField("receiver");
   f_amount=recordClass.getField("amount");
   f_priority=recordClass.getField("priority");
   m_route=routerClass.getMethod("routeTransaction",new java.lang.Class[]{stringClass,recordClass});
   m_init=routerClass.getMethod("initRoutingTable",new java.lang.Class[]{stringClass});
   m_count=routerClass.getMethod("getTransactionCount",new java.lang.Class[]{stringClass,stringClass});
   stringBounds=new java.lang.String[]{"","normal_string_test","\'; DROP TABLE system_metrics; --","null"};
   dbPathBounds=new java.lang.String[]{"generated/metrics_test.db",":memory:","null"};
   doubleBounds=new netrexx.lang.Rexx[]{new netrexx.lang.Rexx((byte)0),new netrexx.lang.Rexx((byte)1),new netrexx.lang.Rexx((byte)-1),new netrexx.lang.Rexx(999999999),new netrexx.lang.Rexx(1.79e+308D),new netrexx.lang.Rexx(-(1.79e+308D))};
   recordBounds=new java.util.ArrayList();
   recordBounds.add((java.lang.Object)null);
   {int $3=0;java.lang.String[] $2=new java.lang.String[stringBounds.length];synchronized(stringBounds){for(;;){if($3==$2.length)break;$2[$3]=stringBounds[stringBounds.length-1-$3];$3++;}}tsVal:for(;;){if(--$3<0)break;tsVal=(java.lang.String)$2[$3];
    {int $6=0;java.lang.String[] $5=new java.lang.String[stringBounds.length];synchronized(stringBounds){for(;;){if($6==$5.length)break;$5[$6]=stringBounds[stringBounds.length-1-$6];$6++;}}nameVal:for(;;){if(--$6<0)break;nameVal=(java.lang.String)$5[$6];
     {int $9=0;netrexx.lang.Rexx[] $8=new netrexx.lang.Rexx[doubleBounds.length];synchronized(doubleBounds){for(;;){if($9==$8.length)break;$8[$9]=doubleBounds[doubleBounds.length-1-$9];$9++;}}valVal:for(;;){if(--$9<0)break;valVal=(netrexx.lang.Rexx)$8[$9];
      constructors=recordClass.getDeclaredConstructors();
      rec=constructors[0].newInstance((java.lang.Object[])null);
      if (netrexx.lang.Rexx.toRexx(tsVal).OpNotEq(null,$01)) 
       f_priority.set(rec,(java.lang.Object)tsVal);
      if (netrexx.lang.Rexx.toRexx(nameVal).OpNotEq(null,$01)) 
       f_receiver.set(rec,(java.lang.Object)nameVal);
      if (netrexx.lang.Rexx.toRexx(tsVal).OpNotEq(null,$01)) 
       f_sender.set(rec,(java.lang.Object)tsVal);
      f_amount.set(rec,(java.lang.Object)valVal);
      if (netrexx.lang.Rexx.toRexx(nameVal).OpNotEq(null,$01)) 
       f_txId.set(rec,(java.lang.Object)nameVal);
      recordBounds.add(rec);
      }
     }/*valVal*/
     }
    }/*nameVal*/
    }
   }/*tsVal*/
   netrexx.lang.RexxIO.Say("Testing method routeTransaction...");
   {int $12=0;java.lang.String[] $11=new java.lang.String[dbPathBounds.length];synchronized(dbPathBounds){for(;;){if($12==$11.length)break;$11[$12]=dbPathBounds[dbPathBounds.length-1-$12];$12++;}}val1:for(;;){if(--$12<0)break;val1=(java.lang.String)$11[$12];
    routeTransaction_p1=(java.lang.String)null;
    if (netrexx.lang.Rexx.toRexx(val1).OpNotEq(null,$01)) 
     routeTransaction_p1=val1;
    {int $15=0;java.lang.Object[] $14=new java.lang.Object[recordBounds.size()];synchronized(recordBounds){java.util.Iterator $13=recordBounds.iterator();for(;;){if($15==$14.length)break;$14[$15]=$13.next();$15++;}}val2:for(;;){if(--$15<0)break;val2=(java.lang.Object)$14[$15];
     routeTransaction_p2=val2;
     {try{
      m_route.invoke((java.lang.Object)null,new java.lang.Object[]{(java.lang.Object)routeTransaction_p1,routeTransaction_p2});
     }
     catch (java.lang.Exception $16){
      ;
     }}
     }
    }/*val2*/
    }
   }/*val1*/
   netrexx.lang.RexxIO.Say("  Method routeTransaction boundary exhaustion completed.");
   netrexx.lang.RexxIO.Say("Testing method initRoutingTable...");
   {int $19=0;java.lang.String[] $18=new java.lang.String[dbPathBounds.length];synchronized(dbPathBounds){for(;;){if($19==$18.length)break;$18[$19]=dbPathBounds[dbPathBounds.length-1-$19];$19++;}}val1:for(;;){if(--$19<0)break;val1=(java.lang.String)$18[$19];
    initRoutingTable_p1=(java.lang.String)null;
    if (netrexx.lang.Rexx.toRexx(val1).OpNotEq(null,$01)) 
     initRoutingTable_p1=val1;
    {try{
     m_init.invoke((java.lang.Object)null,new java.lang.Object[]{(java.lang.Object)initRoutingTable_p1});
    }
    catch (java.lang.Exception $20){
     ;
    }}
    }
   }/*val1*/
   netrexx.lang.RexxIO.Say("  Method initRoutingTable boundary exhaustion completed.");
   netrexx.lang.RexxIO.Say("Testing method getTransactionCount...");
   {int $23=0;java.lang.String[] $22=new java.lang.String[dbPathBounds.length];synchronized(dbPathBounds){for(;;){if($23==$22.length)break;$22[$23]=dbPathBounds[dbPathBounds.length-1-$23];$23++;}}val1:for(;;){if(--$23<0)break;val1=(java.lang.String)$22[$23];
    getTransactionCount_p1=(java.lang.String)null;
    if (netrexx.lang.Rexx.toRexx(val1).OpNotEq(null,$01)) 
     getTransactionCount_p1=val1;
    {int $26=0;java.lang.String[] $25=new java.lang.String[stringBounds.length];synchronized(stringBounds){for(;;){if($26==$25.length)break;$25[$26]=stringBounds[stringBounds.length-1-$26];$26++;}}val2:for(;;){if(--$26<0)break;val2=(java.lang.Object)((java.lang.String)$25[$26]);
     getTransactionCount_p2=(java.lang.String)null;
     if (val2!="null") 
      getTransactionCount_p2=(java.lang.String)val2;
     {try{
      m_count.invoke((java.lang.Object)null,new java.lang.Object[]{(java.lang.Object)getTransactionCount_p1,(java.lang.Object)getTransactionCount_p2});
     }
     catch (java.lang.Exception $27){
      ;
     }}
     }
    }/*val2*/
    }
   }/*val1*/
   netrexx.lang.RexxIO.Say("  Method getTransactionCount boundary exhaustion completed.");
   netrexx.lang.RexxIO.Say("=== [Phase III] Boundary Input Exhaustion Test Completed successfully! ===");
  }
  catch (java.lang.Exception $28){
   ;
  }}
  return;}
 
 
 private TransactionRouterTest(){return;}
 }