/* Generated from 'TransactionRouter.nrx' 15 Jun 2026 00:33:46 [v5.10] */
/* Options: Annotations Binary Decimal Format Implicituses Java Logo Replace Sourcedir Trace2 Verbose3 */
package com.factory.routing;
import java.sql.Connection;
import java.sql.ResultSet;
import java.sql.Statement;
import java.sql.PreparedStatement;
import java.sql.SQLException;
import java.sql.DriverManager;


class TransactionRecordDummy{
 private static final java.lang.String $0="TransactionRouter.nrx";
 
 public TransactionRecordDummy(){return;}
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
 private static final java.lang.String $0="TransactionRouter.nrx";
 
 @SuppressWarnings("unchecked") 
 
 public static void initRoutingTable(java.lang.String dbPath){
  java.sql.Connection conn;
  java.sql.Statement stmt;
  java.lang.Exception ex=null;
  conn=(java.sql.Connection)null;
  stmt=(java.sql.Statement)null;
  {try{
   java.lang.Class.forName("org.sqlite.JDBC");
   conn=DriverManager.getConnection("jdbc:sqlite:"+dbPath);
   stmt=conn.createStatement();
   stmt.executeUpdate("CREATE TABLE IF NOT EXISTS routing_rules (min_amount REAL, priority TEXT, channel TEXT, PRIMARY KEY (min_amount, priority));");
   stmt.executeUpdate("CREATE TABLE IF NOT EXISTS transaction_log (tx_id TEXT PRIMARY KEY, sender TEXT, receiver TEXT, amount REAL, channel TEXT, status TEXT);");
   stmt.executeUpdate("INSERT OR IGNORE INTO routing_rules (min_amount, priority, channel) VALUES (0.0, \'LOW\', \'ACH\');");
   stmt.executeUpdate("INSERT OR IGNORE INTO routing_rules (min_amount, priority, channel) VALUES (0.0, \'HIGH\', \'WIRE\');");
   stmt.executeUpdate("INSERT OR IGNORE INTO routing_rules (min_amount, priority, channel) VALUES (10000.0, \'HIGH\', \'SWIFT\');");
  }
  catch (java.lang.Exception $1){ex=$1;
   ex.printStackTrace();
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
  java.sql.PreparedStatement pstmtSelect;
  java.sql.PreparedStatement pstmtInsert;
  java.sql.ResultSet rs;
  java.lang.String channel;
  java.lang.String status;
  java.lang.String txIdVal=null;
  java.lang.String senderVal=null;
  java.lang.String receiverVal=null;
  double amountVal=0;
  java.lang.String priorityVal=null;
  java.lang.Exception ex=null;
  conn=(java.sql.Connection)null;
  pstmtSelect=(java.sql.PreparedStatement)null;
  pstmtInsert=(java.sql.PreparedStatement)null;
  rs=(java.sql.ResultSet)null;
  channel="NONE";
  status="FAILED";
  if (record!=null) 
   {
    {try{
     java.lang.Class.forName("org.sqlite.JDBC");
     conn=DriverManager.getConnection("jdbc:sqlite:"+dbPath);
     txIdVal="";
     if (record.txId!=null) 
      txIdVal=record.txId;
     else 
      txIdVal=java.util.UUID.randomUUID().toString();
     senderVal="";
     if (record.sender!=null) 
      senderVal=record.sender;
     receiverVal="";
     if (record.receiver!=null) 
      receiverVal=record.receiver;
     amountVal=(double)0.0D;
     if (record.amount!=null) 
      amountVal=record.amount.todouble();
     priorityVal="";
     if (record.priority!=null) 
      priorityVal=record.priority;
     pstmtSelect=conn.prepareStatement("SELECT channel FROM routing_rules WHERE priority = ? AND min_amount <= ? ORDER BY min_amount DESC LIMIT 1;");
     pstmtSelect.setString(1,priorityVal);
     pstmtSelect.setDouble(2,amountVal);
     rs=pstmtSelect.executeQuery();
     if (rs.next()) 
      {
       channel=rs.getString(1);
       status="COMPLETED";
      }
     pstmtInsert=conn.prepareStatement("INSERT OR REPLACE INTO transaction_log (tx_id, sender, receiver, amount, channel, status) VALUES (?, ?, ?, ?, ?, ?);");
     pstmtInsert.setString(1,txIdVal);
     pstmtInsert.setString(2,senderVal);
     pstmtInsert.setString(3,receiverVal);
     pstmtInsert.setDouble(4,amountVal);
     pstmtInsert.setString(5,channel);
     pstmtInsert.setString(6,status);
     pstmtInsert.executeUpdate();
    }
    catch (java.lang.Exception $3){ex=$3;
     ex.printStackTrace();
    }
    finally{
     {try{
      if (rs!=null) 
       rs.close();
      if (pstmtSelect!=null) 
       pstmtSelect.close();
      if (pstmtInsert!=null) 
       pstmtInsert.close();
      if (conn!=null) 
       conn.close();
     }
     catch (java.sql.SQLException $4){
      ;
     }}
    }}
   }
  return channel;
  }
 
 
 @SuppressWarnings("unchecked") 
 
 public static netrexx.lang.Rexx getTransactionCount(java.lang.String dbPath,java.lang.String status){
  java.sql.Connection conn;
  java.sql.PreparedStatement pstmt;
  java.sql.ResultSet rs;
  netrexx.lang.Rexx count;
  java.lang.String statusVal=null;
  java.lang.Exception ex=null;
  conn=(java.sql.Connection)null;
  pstmt=(java.sql.PreparedStatement)null;
  rs=(java.sql.ResultSet)null;
  count=new netrexx.lang.Rexx(0);
  {try{
   java.lang.Class.forName("org.sqlite.JDBC");
   conn=DriverManager.getConnection("jdbc:sqlite:"+dbPath);
   pstmt=conn.prepareStatement("SELECT COUNT(*) FROM transaction_log WHERE status = ?;");
   statusVal="";
   if (status!=null) 
    statusVal=status;
   pstmt.setString(1,statusVal);
   rs=pstmt.executeQuery();
   if (rs.next()) 
    {
     count=new netrexx.lang.Rexx(rs.getInt(1));
    }
  }
  catch (java.lang.Exception $5){ex=$5;
   ex.printStackTrace();
  }
  finally{
   {try{
    if (rs!=null) 
     rs.close();
    if (pstmt!=null) 
     pstmt.close();
    if (conn!=null) 
     conn.close();
   }
   catch (java.sql.SQLException $6){
    ;
   }}
  }}
  return count;
  }
 
 
 @SuppressWarnings("unchecked") 
 
 public static void main(java.lang.String args[]){
  java.lang.String dbFile;
  com.factory.routing.TransactionRecord rec1;
  com.factory.routing.TransactionRecord rec2;
  com.factory.routing.TransactionRecord rec3;
  java.lang.String c1;
  java.lang.String c2;
  java.lang.String c3;
  netrexx.lang.Rexx completedCount;
  dbFile="generated/routing.db";
  initRoutingTable(dbFile);
  rec1=new com.factory.routing.TransactionRecord();
  rec1.txId="tx1";
  rec1.sender="Alice";
  rec1.receiver="Bob";
  rec1.amount=new netrexx.lang.Rexx((short)500);
  rec1.priority="LOW";
  rec2=new com.factory.routing.TransactionRecord();
  rec2.txId="tx2";
  rec2.sender="Charlie";
  rec2.receiver="Dave";
  rec2.amount=new netrexx.lang.Rexx((short)1500);
  rec2.priority="HIGH";
  rec3=new com.factory.routing.TransactionRecord();
  rec3.txId="tx3";
  rec3.sender="Eve";
  rec3.receiver="Frank";
  rec3.amount=new netrexx.lang.Rexx((short)12000);
  rec3.priority="HIGH";
  c1=routeTransaction(dbFile,rec1);
  c2=routeTransaction(dbFile,rec2);
  c3=routeTransaction(dbFile,rec3);
  netrexx.lang.RexxIO.Say("tx1 routed to: "+c1);
  netrexx.lang.RexxIO.Say("tx2 routed to: "+c2);
  netrexx.lang.RexxIO.Say("tx3 routed to: "+c3);
  completedCount=getTransactionCount(dbFile,"COMPLETED");
  netrexx.lang.RexxIO.Say(netrexx.lang.Rexx.toRexx("Completed transactions count: ").OpCc(null,completedCount));
  return;}
 
 
 private TransactionRouter(){return;}
 }