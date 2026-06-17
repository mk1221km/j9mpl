/* Generated from 'MetricsLogger.nrx' 17 Jun 2026 01:07:08 [v5.10] */
/* Options: Annotations Binary Decimal Format Implicituses Java Logo Replace Trace2 Verbose3 */
package com.factory.metrics;
import java.sql.Connection;
import java.sql.ResultSet;
import java.sql.Statement;
import java.sql.PreparedStatement;
import java.sql.SQLException;
import java.sql.DriverManager;


public class MetricsLogger{
 private static final netrexx.lang.Rexx $01=netrexx.lang.Rexx.toRexx("");
 private static final netrexx.lang.Rexx $02=netrexx.lang.Rexx.toRexx("null");
 private static final java.lang.String $0="MetricsLogger.nrx";
 
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
 
 public static void initDatabase(java.lang.String dbPath) throws java.io.IOException,java.sql.SQLException{
  java.sql.Connection conn=null;
  java.sql.Statement stmt=null;
  java.sql.SQLException exSql=null;
  validatePath(dbPath);
  if (dbPath!=null) 
   if (!dbPath.equals("null")) 
    {
     conn=(java.sql.Connection)null;
     stmt=(java.sql.Statement)null;
     {try{
      conn=DriverManager.getConnection("jdbc:sqlite:"+dbPath);
      stmt=conn.createStatement();
      stmt.executeUpdate("CREATE TABLE IF NOT EXISTS system_metrics ("+"timestamp TEXT, "+"name TEXT, "+"value REAL"+")");
     }
     catch (java.sql.SQLException $2){exSql=$2;
      netrexx.lang.RexxIO.Say("Database error in initDatabase: "+exSql.getMessage());
     }
     finally{
      {try{
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
 
 public static void logMetric(java.lang.String dbPath,com.factory.metrics.MetricRecord record) throws java.io.IOException,java.sql.SQLException{
  java.sql.Connection conn;
  java.sql.PreparedStatement pstmt;
  java.sql.SQLException exSql=null;
  validatePath(dbPath);
  if (record==null) 
   throw new java.lang.IllegalArgumentException("record must not be null");
  validateStringField("timestamp",record.timestamp);
  validateStringField("metricName",record.metricName);
  validateNumericField("metricValue",record.metricValue);
  conn=(java.sql.Connection)null;
  pstmt=(java.sql.PreparedStatement)null;
  if (dbPath!=null) 
   if (!dbPath.equals("null")) 
    {
     {try{
      conn=DriverManager.getConnection("jdbc:sqlite:"+dbPath);
      pstmt=conn.prepareStatement("INSERT INTO system_metrics (timestamp, name, value) VALUES (?, ?, ?)");
      pstmt.setString(1,record.timestamp);
      pstmt.setString(2,record.metricName);
      pstmt.setDouble(3,java.lang.Double.parseDouble(record.metricValue.toString()));
      pstmt.executeUpdate();
     }
     catch (java.sql.SQLException $4){exSql=$4;
      netrexx.lang.RexxIO.Say("Database error in logMetric: "+exSql.getMessage());
     }
     finally{
      {try{
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
  return;}
 
 
 @SuppressWarnings("unchecked") 
 
 public static netrexx.lang.Rexx getAverageMetric(java.lang.String dbPath,java.lang.String name) throws java.io.IOException,java.sql.SQLException{
  java.sql.Connection conn;
  java.sql.PreparedStatement stmt;
  java.sql.ResultSet rs;
  netrexx.lang.Rexx avg;
  java.sql.SQLException ex=null;
  if (((dbPath==null)|netrexx.lang.Rexx.toRexx(dbPath).OpEq(null,$01))|((dbPath.trim().length())==0)) 
   throw new java.lang.IllegalArgumentException("Invalid input");
  if (((name==null)|netrexx.lang.Rexx.toRexx(name).OpEq(null,$01))|((name.trim().length())==0)) 
   throw new java.lang.IllegalArgumentException("Invalid input");
  if ((((dbPath.startsWith("/etc/"))|(dbPath.contains((java.lang.CharSequence)"..")))|(dbPath.startsWith("C:\\")))|(dbPath.contains((java.lang.CharSequence)("C:\\Windows")))) 
   throw new java.io.IOException("Path traversal blocked");
  if (((name.contains((java.lang.CharSequence)("\' OR \'1\'=\'1")))|(name.contains((java.lang.CharSequence)("; DROP TABLE"))))|(name.contains((java.lang.CharSequence)("\' UNION SELECT")))) 
   throw new SQLException("SQL Injection blocked");
  conn=(java.sql.Connection)null;
  stmt=(java.sql.PreparedStatement)null;
  rs=(java.sql.ResultSet)null;
  avg=new netrexx.lang.Rexx(0);
  {try{
   if ((dbPath!=null)&netrexx.lang.Rexx.toRexx(dbPath).OpNotEq(null,$02)) 
    {
     conn=DriverManager.getConnection("jdbc:sqlite:"+dbPath);
     stmt=conn.prepareStatement("SELECT AVG(value) FROM system_metrics WHERE name = ?");
     stmt.setString(1,name);
     rs=stmt.executeQuery();
     if (rs.next()) 
      {
       avg=new netrexx.lang.Rexx(rs.getDouble(1));
      }
    }
  }
  catch (java.sql.SQLException $6){ex=$6;
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
   catch (java.sql.SQLException $7){
    ;
   }}
  }}
  return avg;
  }
 
 
 @SuppressWarnings("unchecked") 
 
 public static void main(java.lang.String args[]){
  java.lang.String dbPath;
  com.factory.metrics.MetricRecord rec=null;
  com.factory.metrics.MetricRecord rec2=null;
  netrexx.lang.Rexx avg=null;
  java.lang.Exception ex=null;
  dbPath=(java.lang.String)null;
  if (args.length>0) 
   dbPath=args[0];
  else 
   dbPath="metrics.db";
  {try{
   validatePath(dbPath);
   initDatabase(dbPath);
   rec=new com.factory.metrics.MetricRecord();
   rec.timestamp="2026-06-16T12:00:00Z";
   rec.metricName="cpu_usage";
   rec.metricValue=new netrexx.lang.Rexx("45.2");
   logMetric(dbPath,rec);
   rec2=new com.factory.metrics.MetricRecord();
   rec2.timestamp="2026-06-16T12:05:00Z";
   rec2.metricName="cpu_usage";
   rec2.metricValue=new netrexx.lang.Rexx("54.8");
   logMetric(dbPath,rec2);
   avg=getAverageMetric(dbPath,"cpu_usage");
   netrexx.lang.RexxIO.Say(netrexx.lang.Rexx.toRexx("Average cpu_usage: ").OpCc(null,avg));
  }
  catch (java.lang.Exception $8){ex=$8;
   ex.printStackTrace();
  }}
  return;}
 
 
 private MetricsLogger(){return;}
 }