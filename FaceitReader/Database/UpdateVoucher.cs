using MySql.Data.MySqlClient;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace FaceitReader.Database
{
    internal class UpdateVoucher
    {
        public static void updateVoucher(string teamId, List<string> vouchers)
        {
            string queryU;
            string queryI;

            List<Int32> voucherIds = new List<Int32>();
            MySqlConnection conn = Connection.openConnection();
            MySqlCommand cmd;
            MySqlDataReader reader;

            foreach (var voucher in vouchers)
            {
                queryU = "Select * from vouchers WHERE voucher = '" + voucher + "'";
                cmd = new MySqlCommand(queryU, conn);
                reader = cmd.ExecuteReader();
                reader.Read();
                if(reader.GetString(4) == "created")
                {
                    voucherIds.Add(reader.GetInt32(0));
                }
                reader.Close();
            }

            foreach (var voucher in vouchers)
            {
                queryU = "UPDATE vouchers SET status = 'sendet', team_id = '" + teamId + "', sendet_at = '" + DateTime.Now.ToString("yyyy-MM-dd hh:mm:ss") + "' WHERE voucher = '" + voucher + "' and status = 'created'";
                cmd = new MySqlCommand(queryU, conn);
                cmd.ExecuteNonQuery();
            }

            foreach (var voucher in voucherIds)
            {
                queryI = "INSERT INTO vouchers_user(vouchers_id,user_id) VALUE('"+voucher+"','"+Globals.userId+"')";
                cmd = new MySqlCommand(queryI, conn);
                cmd.ExecuteNonQuery();
            }
            conn.Close();
        }
        public static void reportVoucher(List<string> vouchers)
        {
            string queryU;
            string queryI;
            List<Int32> voucherIds = new List<Int32>();
            MySqlConnection conn = Connection.openConnection();
            MySqlCommand cmd;
            MySqlDataReader reader;

            foreach (var voucher in vouchers)
            {
                queryU = "Select * from vouchers WHERE voucher = '" + voucher + "'";
                cmd = new MySqlCommand(queryU, conn);
                reader = cmd.ExecuteReader();
                while (reader.Read())
                {
                    if (reader.GetString(4) == "sendet")
                    {
                        voucherIds.Add(reader.GetInt32(0));
                    }
                }
                reader.Close();
            }

            foreach (var voucher in vouchers)
            {
                queryU = "UPDATE vouchers SET status = 'not_working' WHERE voucher = '" + voucher + "'";
                cmd = new MySqlCommand(queryU, conn);
                cmd.ExecuteNonQuery();
            }

            foreach (var voucher in voucherIds)
            {
                queryI = "INSERT INTO vouchers_user(vouchers_id,user_id) VALUE('" + voucher + "','" + Globals.userId + "')";
                cmd = new MySqlCommand(queryI, conn);
                cmd.ExecuteNonQuery();
            }
            conn.Close();
        }
        public static void changeVoucher(string voucherOld, string voucher)
        {
            string query = "UPDATE vouchers SET voucher = '" + voucher + "', status = 'changed' WHERE voucher = '" + voucherOld + "'";
            MySqlConnection conn = Connection.openConnection();
            MySqlCommand cmd;

            cmd = new MySqlCommand(query, conn);
            cmd.ExecuteNonQuery();
            conn.Close();
        }
    }
}
