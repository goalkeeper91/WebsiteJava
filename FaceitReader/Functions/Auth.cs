using FaceitReader.Database;
using MySql.Data.MySqlClient;
using System;
using System.Collections.Generic;
using System.Data.SqlClient;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Windows.Forms;

namespace FaceitReader.Functions
{
    internal class Auth
    {
        private static MySqlConnection cnn;

        public static bool userAuth(string username, string password)
        {
            bool auth = false;
            string sql = "Select * from user where username='"+username+"' and password='"+password+"'";
            MySqlCommand command;
            MySqlDataReader reader;
            try
            {
                int a;
                cnn = Connection.openConnection();
                command = new MySqlCommand(sql, cnn);
                reader = command.ExecuteReader();
                if(reader.HasRows)
                {
                    reader.Read();
                    auth = true;
                    if (Convert.ToInt32(reader[4]) == 1)
                    {
                        Globals.admin = true;
                    }
                }
                reader.Close();
                cnn.Close();
            }
            catch (System.Exception)
            {
                MessageBox.Show("Can not open Connection !");
            }
            
            return auth;
        }
    }
}
