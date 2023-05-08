using MySql.Data.MySqlClient;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace FaceitReader.Database
{
    internal class UpdateUser
    {
        public static string updateUser(int id,string update)
        {
            string success = "Erfolgreich";
            string query = "UPDATE user SET " + update + "WHERE id = " + id +";";
            MySqlConnection conn = Connection.openConnection();
            MySqlCommand cmd;

            try
            {
                cmd = new MySqlCommand(query, conn);
                cmd.ExecuteNonQuery();
            }
            catch
            {
                success = "Fehler";
            }

            conn.Close();

            return success;
        }
    }
}
