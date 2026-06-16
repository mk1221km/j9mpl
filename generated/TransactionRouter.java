/* Generated from 'TransactionRouter.nrx' 16 Jun 2026 23:10:42 [v5.10] */
/* Options: Annotations Binary Decimal Format Implicituses Java Logo Replace Trace2 Verbose3 */
package com.factory.routing;
import java.sql.Connection;
import java.sql.ResultSet;
import java.sql.Statement;
import java.sql.PreparedStatement;
import java.sql.SQLException;
import java.sql.DriverManager;


public class TransactionRouter{
 private static final java.lang.String $0="TransactionRouter.nrx";
 
 @SuppressWarnings("unchecked") 
 
 private static void validatePath(java.lang.String path) throws java.io.IOException,java.sql.SQLException{
  if (path==null) 
   throw new java.lang.IllegalArgumentException("path must not be null");
  if ((path.trim().length())==0) 
   throw new java.lang.IllegalArgumentException("path must not be empty");
  if (((((path.indexOf(".."))>=0)|(path.startsWith("/etc/")))|((path.indexOf("C:\\Windows"))>=0))|(path.startsWith("C:\\"))) 
   throw new java.io.IOException("Path traversal blocked");
  if ((((path.indexOf("\' OR \'1\'=\'1"))>=0)|((path.indexOf("; DROP TABLE"))>=0))|((path.indexOf("\' UNION SELECT"))>=0)) 
   throw new SQLException("SQL Injection blocked");
  return;}
 
 
 @SuppressWarnings("unchecked") 
 
 private static void validateStringField(java.lang.String fieldName,java.lang.String value) throws java.io.IOException,java.sql.SQLException{
  if (value==null) 
   throw new java.lang.IllegalArgumentException(fieldName+" must not be null");
  if ((value.trim().length())==0) 
   throw new java.lang.IllegalArgumentException(fieldName+" must not be empty");
  if (((((value.indexOf(".."))>=0)|(value.startsWith("/etc/")))|((value.indexOf("C:\\Windows"))>=0))|(value.startsWith("C:\\"))) 
   throw new java.io.IOException("Path traversal blocked in "+fieldName);
  if ((((value.indexOf("\' OR \'1\'=\'1"))>=0)|((value.indexOf("; DROP TABLE"))>=0))|((value.indexOf("\' UNION SELECT"))>=0)) 
   throw new SQLException("SQL Injection blocked in "+fieldName);
  return;}
 
 
 @SuppressWarnings("unchecked") 
 
 private static void validateNumericField(java.lang.String fieldName,netrexx.lang.Rexx value){
  java.lang.String valStr;
  double dVal=0;
  java.lang.NumberFormatException exNum=null;
  if (value==null) 
   throw new java.lang.NumberFormatException(fieldName+" must not be null");
  valStr=value.toString();
  if ((valStr.trim().length())==0) 
   throw new java.lang.NumberFormatException(fieldName+" must not be empty");
  {try{
   dVal=java.lang.Double.parseDouble(valStr);
   if ((java.lang.Double.isNaN(dVal))|(java.lang.Double.isInfinite(dVal))) 
    throw new java.lang.NumberFormatException("Invalid double");
  }
  catch (java.lang.NumberFormatException $1){exNum=$1;
   throw new java.lang.NumberFormatException("Invalid number format in "+fieldName);
  }}
  return;}
 
 
 @SuppressWarnings("unchecked") 
 
 public static void initRoutingTable(java.lang.String dbPath) throws java.io.IOException,java.sql.SQLException{
  java.sql.Connection conn=null;
  java.sql.Statement stmt=null;
  java.sql.PreparedStatement pstmt=null;
  java.sql.SQLException exSql=null;
  validatePath(dbPath);
  if (dbPath!=null) 
   if (!dbPath.equals("null")) 
    {
     conn=(java.sql.Connection)null;
     stmt=(java.sql.Statement)null;
     pstmt=(java.sql.PreparedStatement)null;
     {try{
      conn=DriverManager.getConnection("jdbc:sqlite:"+dbPath);
      stmt=conn.createStatement();
      stmt.executeUpdate("CREATE TABLE IF NOT EXISTS routing_rules ("+"min_amount REAL, "+"priority TEXT, "+"channel TEXT, "+"PRIMARY KEY (min_amount, priority)"+")");
      stmt.executeUpdate("CREATE TABLE IF NOT EXISTS transaction_log ("+"tx_id TEXT PRIMARY KEY, "+"sender TEXT, "+"receiver TEXT, "+"amount REAL, "+"channel TEXT, "+"status TEXT"+")");
      pstmt=conn.prepareStatement("INSERT OR IGNORE INTO routing_rules "+"(min_amount, priority, channel) VALUES (?, ?, ?)");
      pstmt.setDouble(1,0.0D);
      pstmt.setString(2,"low");
      pstmt.setString(3,"bank");
      pstmt.executeUpdate();
      pstmt.setDouble(1,1000.0D);
      pstmt.setString(2,"medium");
      pstmt.setString(3,"card");
      pstmt.executeUpdate();
      pstmt.setDouble(1,5000.0D);
      pstmt.setString(2,"high");
      pstmt.setString(3,"wire");
      pstmt.executeUpdate();
     }
     catch (java.sql.SQLException $2){exSql=$2;
      netrexx.lang.RexxIO.Say("Database error in initRoutingTable: "+exSql.getMessage());
     }
     finally{
      {try{
       if (pstmt!=null) 
        pstmt.close();
       if (stmt!=null) 
        stmt.close();
       if (conn!=null) 
        conn.close();
      }
      catch (java.sql.SQLException $3){
       ;
      }}
     }}
    }
  return;}
 
 
 @SuppressWarnings("unchecked") 
 
 public static java.lang.String routeTransaction(java.lang.String dbPath,com.factory.routing.TransactionRecord record) throws java.io.IOException,java.sql.SQLException{
  java.sql.Connection conn;
  java.sql.PreparedStatement stmt;
  java.sql.PreparedStatement pstmt;
  java.sql.ResultSet rs;
  java.lang.String channel;
  java.sql.SQLException exSql=null;
  validatePath(dbPath);
  if (record==null) 
   throw new java.lang.IllegalArgumentException("record must not be null");
  validateStringField("txId",record.txId);
  validateStringField("sender",record.sender);
  validateStringField("receiver",record.receiver);
  validateStringField("priority",record.priority);
  validateNumericField("amount",record.amount);
  conn=(java.sql.Connection)null;
  stmt=(java.sql.PreparedStatement)null;
  pstmt=(java.sql.PreparedStatement)null;
  rs=(java.sql.ResultSet)null;
  channel=(java.lang.String)null;
  if (dbPath!=null) 
   if (!dbPath.equals("null")) 
    {
     {try{
      conn=DriverManager.getConnection("jdbc:sqlite:"+dbPath);
      stmt=conn.prepareStatement("SELECT channel FROM routing_rules WHERE min_amount <= ? AND LOWER(priority) = LOWER(?) ORDER BY min_amount DESC LIMIT 1");
      stmt.setDouble(1,java.lang.Double.parseDouble(record.amount.toString()));
      stmt.setString(2,record.priority);
      rs=stmt.executeQuery();
      if (rs.next()) 
       {
        channel=rs.getString(1);
       }
      else 
       {
        channel="wire";
       }
      pstmt=conn.prepareStatement("INSERT INTO transaction_log (tx_id, sender, receiver, amount, channel, status) VALUES (?, ?, ?, ?, ?, ?)");
      pstmt.setString(1,record.txId);
      pstmt.setString(2,record.sender);
      pstmt.setString(3,record.receiver);
      pstmt.setDouble(4,java.lang.Double.parseDouble(record.amount.toString()));
      pstmt.setString(5,channel);
      pstmt.setString(6,"queued");
      pstmt.executeUpdate();
     }
     catch (java.sql.SQLException $4){exSql=$4;
      netrexx.lang.RexxIO.Say("Database error in routeTransaction: "+exSql.getMessage());
     }
     finally{
      {try{
       if (rs!=null) 
        rs.close();
       if (stmt!=null) 
        stmt.close();
       if (pstmt!=null) 
        pstmt.close();
       if (conn!=null) 
        conn.close();
      }
      catch (java.sql.SQLException $5){
       ;
      }}
     }}
    }
  return channel;
  }
 
 
 @SuppressWarnings("unchecked") 
 
 public static netrexx.lang.Rexx getTransactionCount(java.lang.String dbPath,java.lang.String status) throws java.io.IOException,java.sql.SQLException{
  java.sql.Connection conn;
  java.sql.PreparedStatement stmt;
  java.sql.ResultSet rs;
  netrexx.lang.Rexx count;
  java.sql.SQLException exSql=null;
  validatePath(dbPath);
  validateStringField("status",status);
  conn=(java.sql.Connection)null;
  stmt=(java.sql.PreparedStatement)null;
  rs=(java.sql.ResultSet)null;
  count=new netrexx.lang.Rexx(0);
  if (dbPath!=null) 
   if (!dbPath.equals("null")) 
    {
     {try{
      conn=DriverManager.getConnection("jdbc:sqlite:"+dbPath);
      if ((status.equals((java.lang.Object)"ALL"))|(status.equals((java.lang.Object)"all"))) 
       {
        stmt=conn.prepareStatement("SELECT COUNT(*) FROM transaction_log");
       }
      else 
       {
        stmt=conn.prepareStatement("SELECT COUNT(*) FROM transaction_log WHERE LOWER(status) = LOWER(?)");
        stmt.setString(1,status);
       }
      rs=stmt.executeQuery();
      if (rs.next()) 
       {
        count=new netrexx.lang.Rexx(rs.getInt(1));
       }
     }
     catch (java.sql.SQLException $6){exSql=$6;
      netrexx.lang.RexxIO.Say("Database error in getTransactionCount: "+exSql.getMessage());
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
      catch (java.sql.SQLException $7){
       ;
      }}
     }}
    }
  return count;
  }
 
 
 @SuppressWarnings("unchecked") 
 
 public static void main(java.lang.String args[]){
  java.lang.String dbPath;
  com.factory.routing.TransactionRecord record1=null;
  com.factory.routing.TransactionRecord record2=null;
  java.lang.String ch1=null;
  java.lang.String ch2=null;
  netrexx.lang.Rexx highCount=null;
  netrexx.lang.Rexx totalCount=null;
  java.lang.Exception ex=null;
  dbPath=(java.lang.String)null;
  if (args.length>0) 
   dbPath=args[0];
  else 
   dbPath="runtime.db";
  {try{
   validatePath(dbPath);
   initRoutingTable(dbPath);
   record1=new com.factory.routing.TransactionRecord();
   record1.txId="TXN001";
   record1.sender="Alice";
   record1.receiver="Bob";
   record1.amount=new netrexx.lang.Rexx("150.50");
   record1.priority="high";
   record2=new com.factory.routing.TransactionRecord();
   record2.txId="TXN002";
   record2.sender="Charlie";
   record2.receiver="Dave";
   record2.amount=new netrexx.lang.Rexx("75.00");
   record2.priority="low";
   ch1=routeTransaction(dbPath,record1);
   ch2=routeTransaction(dbPath,record2);
   netrexx.lang.RexxIO.Say("Route 1 Channel: "+ch1);
   netrexx.lang.RexxIO.Say("Route 2 Channel: "+ch2);
   highCount=getTransactionCount(dbPath,"queued");
   totalCount=getTransactionCount(dbPath,"ALL");
   netrexx.lang.RexxIO.Say(netrexx.lang.Rexx.toRexx("Queued count: ").OpCc(null,highCount));
   netrexx.lang.RexxIO.Say(netrexx.lang.Rexx.toRexx("Total count: ").OpCc(null,totalCount));
  }
  catch (java.lang.Exception $8){ex=$8;
   ex.printStackTrace();
  }}
  return;}
 
 
 private TransactionRouter(){return;}
 }