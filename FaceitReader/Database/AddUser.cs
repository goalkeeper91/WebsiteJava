using FaceitReader.Functions;
using MySql.Data.MySqlClient;

namespace FaceitReader.Database
{
    internal class AddUserToDatabase
    {
        public static string CheckIfUserExist(string username, string password, string email, string firstname, string lastname, int admin)
        {
            string exist = null;
            string query = "Select * from user where username = '" + username + "' and email = '" + email + "'";
            MySqlConnection conn = Connection.openConnection();
            MySqlCommand cmd;
            MySqlDataReader rdr;

            cmd = new MySqlCommand(query, conn);
            rdr = cmd.ExecuteReader();

            if (!rdr.HasRows)
            {
                AddUser(username, password, email, firstname, lastname, admin);
            }
            else
            {
                exist = "User existiert bereits!";
            }

            rdr.Close();
            conn.Close();

            return exist;
        }
        public static void AddUser(string username, string password, string email, string firstname, string lastname, int admin)
        {
            MySqlConnection cnn = Connection.openConnection();
            MySqlCommand command;
            string pw = EncryptString.Encrypt(password);
            string query = "Insert Into user (username,password,email,firstname,lastname,admin) Values('" + username + "','" + pw + "','" + email + "','" + firstname + "','" + lastname + "','" + admin + "')";

            command = new MySqlCommand(query, cnn);
            command.ExecuteNonQuery();

            cnn.Close();
        }
    }
}
