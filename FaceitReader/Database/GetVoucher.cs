using FaceitReader.Classes;
using MySql.Data.MySqlClient;
using System;
using System.Collections.Generic;

namespace FaceitReader.Database
{
    internal class GetVoucher
    {
        public static List<string> getFiveVoucher(string value, string cupname)
        {
            MySqlConnection cnn = Connection.openConnection();
            MySqlCommand command;
            MySqlDataReader reader;
            List<string> vouchers = new List<string>();
            string query = "Select voucher From vouchers where value = " + value + " AND status = 'created' AND cupname = '" + cupname + "' LIMIT 5";
            command = new MySqlCommand(query, cnn);
            reader = command.ExecuteReader();
            if (reader.HasRows)
            {
                while (reader.Read())
                {
                    vouchers.Add(reader.GetString(0));
                }
            }
            return vouchers;
        }

        public static List<Tuple<string,string>> getSendetVoucher(string teamname, string cupname)
        {
            MySqlConnection cnn = Connection.openConnection();
            MySqlCommand command;
            MySqlDataReader reader;
            List<Tuple<string, string>> vouchers = new List<Tuple<string, string>>();
            string query = "Select voucher, status From vouchers where team_id = '" + teamname + "' AND cupname = '" + cupname + "'";
            command = new MySqlCommand(query, cnn);
            reader = command.ExecuteReader();
            if (reader.HasRows)
            {
                while (reader.Read())
                {
                    vouchers.Add(new Tuple<string,string>(reader.GetString(0),reader.GetString(1)));
                }
            }
            Console.WriteLine(vouchers.Count);
            return vouchers;
        }

        public static List<Voucher> getAllVouchersPolska()
        {
            string query = "SELECT * FROM vouchers LEFT JOIN vouchers_user ON vouchers.id = vouchers_user.vouchers_id LEFT JOIN user ON user.id = vouchers_user.user_id WHERE vouchers.cupname like '%POLSKA%';";
            List<Voucher> vouchers = new List<Voucher>();
            MySqlConnection cnn = Connection.openConnection();
            MySqlCommand command;
            MySqlDataReader reader;

            command = new MySqlCommand(query, cnn);
            reader = command.ExecuteReader();
            if (reader.HasRows)
            {
                while (reader.Read())
                {
                    string username;
                    string teamId;
                    string sendetAt;
                    string updatedAt;
                    if (!reader.IsDBNull(14))
                    {
                        username = reader.GetString(14);
                    }
                    else
                    {
                        username = null;
                    }
                    if (!reader.IsDBNull(5))
                    {
                        teamId = reader.GetString(5);
                    }
                    else
                    {
                        teamId = null;
                    }
                    if (!reader.IsDBNull(6))
                    {
                        sendetAt = reader.GetString(6);
                    }
                    else
                    {
                        sendetAt = null;
                    }
                    if (!reader.IsDBNull(8))
                    {
                        updatedAt = reader.GetString(8);
                    }
                    else
                    {
                        updatedAt = null;
                    }
                    vouchers.Add(new Voucher() { Id = reader.GetInt32(0), Cupname = reader.GetString(1), Value = reader.GetString(2), Code = reader.GetString(3), Status = reader.GetString(4), TeamId = teamId, SendetAt = sendetAt, UpdatedFrom = username, Created = reader.GetString(7), Updated = updatedAt });
                }
            }
            cnn.Close();
            
            return vouchers;
        }

        public static List<Voucher> getAllVouchersSk()
        {
            string query = "SELECT * FROM vouchers LEFT JOIN vouchers_user ON vouchers.id = vouchers_user.vouchers_id LEFT JOIN user ON user.id = vouchers_user.user_id WHERE vouchers.cupname like '%SKINBARON%';";
            List<Voucher> vouchers = new List<Voucher>();
            MySqlConnection cnn = Connection.openConnection();
            MySqlCommand command;
            MySqlDataReader reader;

            command = new MySqlCommand(query, cnn);
            reader = command.ExecuteReader();
            if (reader.HasRows)
            {
                while (reader.Read())
                {
                    string username;
                    string teamId;
                    string sendetAt;
                    string updatedAt;
                    if (!reader.IsDBNull(14))
                    {
                        username = reader.GetString(14);
                    }
                    else
                    {
                        username = null;
                    }
                    if (!reader.IsDBNull(5))
                    {
                        teamId = reader.GetString(5);
                    }
                    else
                    {
                        teamId = null;
                    }
                    if (!reader.IsDBNull(6))
                    {
                        sendetAt = reader.GetString(6);
                    }
                    else
                    {
                        sendetAt = null;
                    }
                    if (!reader.IsDBNull(8))
                    {
                        updatedAt = reader.GetString(8);
                    }
                    else
                    {
                        updatedAt = null;
                    }
                    vouchers.Add(new Voucher() { Id = reader.GetInt32(0), Cupname = reader.GetString(1), Value = reader.GetString(2), Code = reader.GetString(3), Status = reader.GetString(4), TeamId = teamId, SendetAt = sendetAt, UpdatedFrom = username, Created = reader.GetString(7), Updated = updatedAt });
                }
            }
            cnn.Close();
            return vouchers;
        }
    }
}
