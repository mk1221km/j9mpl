/* Generated from 'TransactionRouter.nrx' 16 Jun 2026 17:20:45 [v5.10] */
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
  java.sql.PreparedStatement pstmt=null;
  java.sql.SQLException ex=null;
  conn=(java.sql.Connection)null;
  stmt=(java.sql.Statement)null;
  if ((dbPath!=null)&netrexx.lang.Rexx.toRexx(dbPath).OpNotEq(null,$01)) 
   {
    {try{
     conn=DriverManager.getConnection("jdbc:sqlite:"+dbPath);
     stmt=conn.createStatement();
     stmt.executeUpdate("CREATE TABLE IF NOT EXISTS routing_rules ("+"min_amount REAL, "+"priority TEXT, "+"channel TEXT, "+"PRIMARY KEY (min_amount, priority))");
     stmt.executeUpdate("CREATE TABLE IF NOT EXISTS transaction_log ("+"tx_id TEXT PRIMARY KEY, "+"sender TEXT, "+"receiver TEXT, "+"amount REAL, "+"channel TEXT, "+"status TEXT)");
     pstmt=(java.sql.PreparedStatement)null;
     {try{
      pstmt=conn.prepareStatement("INSERT OR IGNORE INTO routing_rules (min_amount, priority, channel) VALUES (?, ?, ?)");
      pstmt.setDouble(1,(double)((byte)0));
      pstmt.setString(2,"low");
      pstmt.setString(3,"default");
      pstmt.executeUpdate();
      pstmt.setDouble(1,(double)((byte)100));
      pstmt.setString(2,"medium");
      pstmt.setString(3,"fast");
      pstmt.executeUpdate();
      pstmt.setDouble(1,(double)((short)1000));
      pstmt.setString(2,"high");
      pstmt.setString(3,"premium");
      pstmt.executeUpdate();
     }
     catch (java.sql.SQLException $1){ex=$1;
      netrexx.lang.RexxIO.Say("Error inserting default routing rules: "+ex.getMessage());
     }
     finally{
      {try{
       if (pstmt!=null) 
        pstmt.close();
      }
      catch (java.sql.SQLException $2){
       ;
      }}
     }}
    }
    catch (java.sql.SQLException $3){ex=$3;
     netrexx.lang.RexxIO.Say("Database initialization error: "+ex.getMessage());
    }
    finally{
     {try{
      if (stmt!=null) 
       stmt.close();
      if (conn!=null) 
       conn.close();
     }
     catch (java.sql.SQLException $4){
      ;
     }}
    }}
   }
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
  channel="UNKNOWN";
  if ((dbPath!=null)&netrexx.lang.Rexx.toRexx(dbPath).OpNotEq(null,$01)) 
   {
    {try{
     conn=DriverManager.getConnection("jdbc:sqlite:"+dbPath);
     stmt=conn.prepareStatement("SELECT channel FROM routing_rules WHERE priority = ? AND min_amount <= ? ORDER BY min_amount DESC LIMIT 1");
     stmt.setString(1,record.priority);
     stmt.setDouble(2,record.amount.todouble());
     rs=stmt.executeQuery();
     if (rs.next()) 
      {
       channel=rs.getString("channel");
      }
     stmt=conn.prepareStatement("INSERT INTO transaction_log (tx_id, sender, receiver, amount, channel, status) VALUES (?, ?, ?, ?, ?, ?)");
     stmt.setString(1,record.txId);
     stmt.setString(2,record.sender);
     stmt.setString(3,record.receiver);
     stmt.setDouble(4,record.amount.todouble());
     stmt.setString(5,channel);
     stmt.setString(6,"ROUTED");
     stmt.executeUpdate();
    }
    catch (java.sql.SQLException $5){ex=$5;
     netrexx.lang.RexxIO.Say("Routing error: "+ex.getMessage());
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
  return channel;
  }
 
 
 @SuppressWarnings("unchecked") 
 
 public static netrexx.lang.Rexx getTransactionCount(java.lang.String dbPath,java.lang.String status){
  java.sql.Connection conn;
  java.sql.PreparedStatement stmt;
  java.sql.ResultSet rs;
  int count;
  java.sql.SQLException ex=null;
  conn=(java.sql.Connection)null;
  stmt=(java.sql.PreparedStatement)null;
  rs=(java.sql.ResultSet)null;
  count=0;
  if ((dbPath!=null)&netrexx.lang.Rexx.toRexx(dbPath).OpNotEq(null,$01)) 
   {
    {try{
     conn=DriverManager.getConnection("jdbc:sqlite:"+dbPath);
     stmt=conn.prepareStatement("SELECT count(*) FROM transaction_log WHERE status = ?");
     stmt.setString(1,status);
     rs=stmt.executeQuery();
     if (rs.next()) 
      {
       count=rs.getInt(1);
      }
    }
    catch (java.sql.SQLException $7){ex=$7;
     netrexx.lang.RexxIO.Say("Database query error: "+ex.getMessage());
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
     catch (java.sql.SQLException $8){
      ;
     }}
    }}
   }
  return new netrexx.lang.Rexx(count);
  }
 
 
 @SuppressWarnings("unchecked") 
 
 public static void main(java.lang.String args[]){
  java.lang.String dbPath;
  java.sql.Connection conn=null;
  java.sql.PreparedStatement stmt=null;
  java.sql.SQLException ex=null;
  com.factory.routing.TransactionRecord rec1=null;
  com.factory.routing.TransactionRecord rec2=null;
  netrexx.lang.Rexx count1=null;
  dbPath=(java.lang.String)null;
  if (args.length>0) 
   dbPath=args[0];
  else 
   dbPath="routing.db";
  if ((dbPath!=null)&netrexx.lang.Rexx.toRexx(dbPath).OpNotEq(null,$01)) 
   {
    com.factory.routing.TransactionRouter.initRoutingTable(dbPath);
    conn=(java.sql.Connection)null;
    stmt=(java.sql.PreparedStatement)null;
    {try{
     conn=DriverManager.getConnection("jdbc:sqlite:"+dbPath);
     stmt=conn.prepareStatement("INSERT OR IGNORE INTO routing_rules (min_amount, priority, channel) VALUES (?, ?, ?)");
     stmt.setDouble(1,0.0D);
     stmt.setString(2,"high");
     stmt.setString(3,"fast");
     stmt.executeUpdate();
     stmt.setDouble(1,100.0D);
     stmt.setString(2,"low");
     stmt.setString(3,"slow");
     stmt.executeUpdate();
    }
    catch (java.sql.SQLException $9){ex=$9;
     netrexx.lang.RexxIO.Say("Error inserting default rules: "+ex.getMessage());
    }
    finally{
     {try{
      if (stmt!=null) 
       stmt.close();
      if (conn!=null) 
       conn.close();
     }
     catch (java.sql.SQLException $10){
      ;
     }}
    }}
    rec1=new com.factory.routing.TransactionRecord();
    rec1.txId="tx001";
    rec1.sender="Alice";
    rec1.receiver="Bob";
    rec1.amount=new netrexx.lang.Rexx("50.0");
    rec1.priority="high";
    com.factory.routing.TransactionRouter.routeTransaction(dbPath,rec1);
    rec2=new com.factory.routing.TransactionRecord();
    rec2.txId="tx002";
    rec2.sender="Charlie";
    rec2.receiver="Dave";
    rec2.amount=new netrexx.lang.Rexx("200.0");
    rec2.priority="low";
    com.factory.routing.TransactionRouter.routeTransaction(dbPath,rec2);
    count1=com.factory.routing.TransactionRouter.getTransactionCount(dbPath,"routed");
    netrexx.lang.RexxIO.Say(netrexx.lang.Rexx.toRexx("Total routed transactions: ").OpCc(null,count1));
   }
  return;}
 
 
 private TransactionRouter(){return;}
 }