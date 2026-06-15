/* Generated from 'TransactionRouter.nrx' 15 Jun 2026 05:00:40 [v5.10] */
/* Options: Annotations Binary Decimal Format Implicituses Java Logo Replace Trace2 Verbose3 */
package com.factory.routing;
import java.sql.Connection;
import java.sql.ResultSet;
import java.sql.Statement;
import java.sql.PreparedStatement;
import java.sql.SQLException;
import java.sql.DriverManager;


class TransactionRouterDummy{
 private static final java.lang.String $0="TransactionRouter.nrx";
 
 public TransactionRouterDummy(){return;}
 }


class TransactionRecord{
 private static final java.lang.String $0="TransactionRouter.nrx";
 /* properties public */
 public java.lang.String txId;
 public java.lang.String sender;
 public java.lang.String receiver;
 public netrexx.lang.Rexx amount;
 public java.lang.String priority;
 
 public TransactionRecord(){return;}
 }


public class TransactionRouter{
 private static final netrexx.lang.Rexx $01=netrexx.lang.Rexx.toRexx("null");
 private static final java.lang.String $0="TransactionRouter.nrx";
 
 @SuppressWarnings("unchecked") 
 
 public static void initRoutingTable(java.lang.String dbPath){
  java.sql.Connection conn;
  java.sql.Statement stmt;
  java.sql.SQLException ex=null;
  conn=(java.sql.Connection)null;
  stmt=(java.sql.Statement)null;
  {try{
   if ((dbPath!=null)&netrexx.lang.Rexx.toRexx(dbPath).OpNotEq(null,$01)) 
    {
     conn=DriverManager.getConnection("jdbc:sqlite:"+dbPath);
     stmt=conn.createStatement();
     stmt.executeUpdate("CREATE TABLE IF NOT EXISTS routing_rules ("+"min_amount REAL, "+"priority TEXT, "+"channel TEXT, "+"PRIMARY KEY (min_amount, priority)"+")");
     stmt.executeUpdate("CREATE TABLE IF NOT EXISTS transaction_log ("+"tx_id TEXT PRIMARY KEY, "+"sender TEXT, "+"receiver TEXT, "+"amount REAL, "+"channel TEXT, "+"status TEXT"+")");
     stmt.executeUpdate("INSERT OR IGNORE INTO routing_rules (min_amount, priority, channel) VALUES (0, \'low\', \'standard\')");
     stmt.executeUpdate("INSERT OR IGNORE INTO routing_rules (min_amount, priority, channel) VALUES (5000, \'medium\', \'premium\')");
     stmt.executeUpdate("INSERT OR IGNORE INTO routing_rules (min_amount, priority, channel) VALUES (10000, \'high\', \'gold\')");
    }
  }
  catch (java.sql.SQLException $1){ex=$1;
   netrexx.lang.RexxIO.Say("SQL error: "+ex.getMessage());
  }
  finally{
   {try{
    if (stmt!=null) 
     stmt.close();
    if (conn!=null) 
     conn.close();
   }
   catch (java.sql.SQLException $2){
    ;
   }}
  }}
  return;}
 
 
 @SuppressWarnings("unchecked") 
 
 public static java.lang.String routeTransaction(java.lang.String dbPath,com.factory.routing.TransactionRecord record){
  java.sql.Connection conn;
  java.sql.PreparedStatement stmt;
  java.sql.ResultSet rs;
  java.lang.String channel;
  java.sql.SQLException ex=null;
  conn=(java.sql.Connection)null;
  stmt=(java.sql.PreparedStatement)null;
  rs=(java.sql.ResultSet)null;
  channel=(java.lang.String)null;
  {try{
   if ((dbPath!=null)&netrexx.lang.Rexx.toRexx(dbPath).OpNotEq(null,$01)) 
    {
     conn=DriverManager.getConnection("jdbc:sqlite:"+dbPath);
     stmt=conn.prepareStatement("SELECT channel FROM routing_rules WHERE min_amount <= ? AND priority = ? ORDER BY min_amount DESC LIMIT 1");
     stmt.setDouble(1,record.amount.todouble());
     stmt.setString(2,record.priority);
     rs=stmt.executeQuery();
     if (rs.next()) 
      {
       channel=rs.getString("channel");
      }
     else 
      {
       channel="ACH";
      }
     stmt=conn.prepareStatement("INSERT INTO transaction_log (tx_id, sender, receiver, amount, channel, status) VALUES (?, ?, ?, ?, ?, ?)");
     stmt.setString(1,record.txId);
     stmt.setString(2,record.sender);
     stmt.setString(3,record.receiver);
     stmt.setDouble(4,record.amount.todouble());
     stmt.setString(5,channel);
     stmt.setString(6,"PENDING");
     stmt.executeUpdate();
    }
   else 
    {
     channel="INVALID_DB";
    }
  }
  catch (java.sql.SQLException $3){ex=$3;
   netrexx.lang.RexxIO.Say("Routing Error: "+ex.getMessage());
   channel="ERROR";
  }
  finally{
   {try{
    if (rs!=null) 
     rs.close();
    if (stmt!=null) 
     stmt.close();
    if (conn!=null) 
     conn.close();
   }
   catch (java.sql.SQLException $4){
    ;
   }}
  }}
  return channel;
  }
 
 
 @SuppressWarnings("unchecked") 
 
 public static netrexx.lang.Rexx getTransactionCount(java.lang.String dbPath,java.lang.String status){
  int count;
  java.sql.Connection conn=null;
  java.sql.PreparedStatement stmt=null;
  java.sql.ResultSet rs=null;
  java.sql.SQLException ex=null;
  count=0;
  if ((dbPath!=null)&netrexx.lang.Rexx.toRexx(dbPath).OpNotEq(null,$01)) 
   {
    conn=(java.sql.Connection)null;
    stmt=(java.sql.PreparedStatement)null;
    rs=(java.sql.ResultSet)null;
    {try{
     conn=DriverManager.getConnection("jdbc:sqlite:"+dbPath);
     stmt=conn.prepareStatement("SELECT count(*) FROM transaction_log WHERE status = ?");
     stmt.setString(1,status);
     rs=stmt.executeQuery();
     {for(;;){if(!(rs.next()))break;
      count=rs.getInt(1);
      }
     }
    }
    catch (java.sql.SQLException $5){ex=$5;
     netrexx.lang.RexxIO.Say("Execution Fault: "+ex.getMessage());
    }
    finally{
     {try{
      if (rs!=null) 
       rs.close();
      if (stmt!=null) 
       stmt.close();
      if (conn!=null) 
       conn.close();
     }
     catch (java.sql.SQLException $6){
      ;
     }}
    }}
   }
  return new netrexx.lang.Rexx(count);
  }
 
 
 @SuppressWarnings("unchecked") 
 
 public static void main(java.lang.String args[]){
  java.lang.String dbPath=null;
  com.factory.routing.TransactionRecord record1=null;
  java.lang.String res1=null;
  com.factory.routing.TransactionRecord record2=null;
  java.lang.String res2=null;
  com.factory.routing.TransactionRecord record3=null;
  java.lang.String res3=null;
  netrexx.lang.Rexx count=null;
  if (args.length>0) 
   dbPath=args[0];
  else 
   dbPath="routing.db";
  if ((dbPath!=null)&netrexx.lang.Rexx.toRexx(dbPath).OpNotEq(null,$01)) 
   {
    initRoutingTable(dbPath);
    record1=new com.factory.routing.TransactionRecord();
    record1.txId="TXN-001";
    record1.sender="Alice";
    record1.receiver="Bob";
    record1.amount=new netrexx.lang.Rexx("100.50");
    record1.priority="low";
    res1=routeTransaction(dbPath,record1);
    netrexx.lang.RexxIO.Say("Routed TXN-001 to channel: "+res1);
    record2=new com.factory.routing.TransactionRecord();
    record2.txId="TXN-002";
    record2.sender="Charlie";
    record2.receiver="Dave";
    record2.amount=new netrexx.lang.Rexx("5500.00");
    record2.priority="medium";
    res2=routeTransaction(dbPath,record2);
    netrexx.lang.RexxIO.Say("Routed TXN-002 to channel: "+res2);
    record3=new com.factory.routing.TransactionRecord();
    record3.txId="TXN-003";
    record3.sender="Eve";
    record3.receiver="Frank";
    record3.amount=new netrexx.lang.Rexx("10500.00");
    record3.priority="high";
    res3=routeTransaction(dbPath,record3);
    netrexx.lang.RexxIO.Say("Routed TXN-003 to channel: "+res3);
    count=getTransactionCount(dbPath,"PENDING");
    netrexx.lang.RexxIO.Say(netrexx.lang.Rexx.toRexx("Total pending transactions: ").OpCc(null,count));
   }
  return;}
 
 
 private TransactionRouter(){return;}
 }