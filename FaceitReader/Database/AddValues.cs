using MySql.Data.MySqlClient;
using Newtonsoft.Json;
using Newtonsoft.Json.Linq;
using System.IO;
using System.Net;

namespace FaceitReader.Database
{
    internal class AddValues
    {
        public static void fillData(string id, int match)
        {
            var Murl = "https://open.faceit.com/data/v4/championships/" + id + "/matches?type=past&offset=0&limit=" + match;
            var httpMRequest = (HttpWebRequest)WebRequest.Create(Murl);


            httpMRequest.ContentType = "application/json; charset=utf-8";
            httpMRequest.Accept = "application/json";
            httpMRequest.Headers["Authorization"] = Globals.Api;
            var httpMResponse = (HttpWebResponse)httpMRequest.GetResponse();
            MySqlConnection cnn = Connection.openConnection();
            MySqlConnection cnn2 = Connection.openConnection();
            MySqlConnection cnn3 = Connection.openConnection();
            MySqlConnection cnn4 = Connection.openConnection();
            MySqlCommand command;
            MySqlDataReader reader1;
            MySqlDataReader reader2;
            MySqlDataReader reader3;
            MySqlDataReader reader4;


            using (var streamReader = new StreamReader(httpMResponse.GetResponseStream()))
            {
                var result = streamReader.ReadToEnd();
                var x = 0;

                JObject json = JsonConvert.DeserializeObject<JObject>(result);

                foreach (JToken item in json["items"])
                {
                    string queryCup = "Select * from cups where cupid ='" + item["competition_id"].ToString() + "'";
                    string queryCupInsert = "INSERT INTO cups (cupid,cupname) VALUES('" + item["competition_id"].ToString() + "','" + item["competition_name"].ToString() + "')";
                    command = new MySqlCommand(queryCup, cnn3);
                    reader3 = command.ExecuteReader();
                    if (!reader3.HasRows)
                    {
                        command = new MySqlCommand(queryCupInsert, cnn);
                        command.ExecuteNonQuery();
                    }
                    reader3.Close();
                    foreach (JToken placement in item["teams"])
                    {
                        foreach (JToken team in placement.Children())
                        {
                            queryCup = "Select * from places_cups where cup_id ='" + item["competition_id"].ToString() + "'";
                            command = new MySqlCommand(queryCup, cnn4);
                            reader4 = command.ExecuteReader();
                            if (!reader4.HasRows)
                            {
                                string queryT = "Select * from teams where team_id ='" + team["faction_id"].ToString() + "'";
                                string queryI = "INSERT INTO teams (team_id,captain_id,teamname) VALUES('" + team["faction_id"].ToString() + "','" + team["leader"].ToString() + "','" + team["name"].ToString() + "')";
                                command = new MySqlCommand(queryT, cnn2);
                                reader2 = command.ExecuteReader();
                                if (!reader2.HasRows)
                                {
                                    if (team["faction_id"].ToString() != "bye")
                                    {
                                        command = new MySqlCommand(queryI, cnn);
                                        command.ExecuteNonQuery();
                                    }
                                }
                                reader2.Close();

                                if (team["faction_id"].ToString() != "bye")
                                {
                                    foreach (JToken player in team["roster"])
                                    {
                                        x++;
                                        string query2 = "Select * from player where player_id ='" + player["player_id"].ToString() + "'";
                                        string query = "";
                                        command = new MySqlCommand(query2, cnn2);
                                        reader1 = command.ExecuteReader();
                                        if (!reader1.HasRows)
                                        {
                                            if (team["faction_id"].ToString() != "bye")
                                            {
                                                query = "INSERT INTO player (team_id,player_id,cup_id,playername,game_player_id) VALUES('" + team["faction_id"].ToString() + "','" + player["player_id"].ToString() + "','" + id + "','" + player["nickname"].ToString() + "','" + player["game_player_id"].ToString() + "')";
                                            }
                                        }
                                        else
                                        {
                                            query = "UPDATE player SET cup_id='" + id + "' WHERE player_id = '" + player["player_id"].ToString() + "'";
                                        }
                                        reader1.Close();
                                        command = new MySqlCommand(query, cnn);
                                        command.ExecuteNonQuery();
                                    }
                                }
                            }
                            else
                            {
                                while (reader4.Read())
                                {
                                    if (reader4["teams_id"] != team["faction_id"])
                                    {
                                        string queryT = "Select * from teams where team_id ='" + team["faction_id"].ToString() + "'";
                                        string queryI = "INSERT INTO teams (team_id,captain_id,teamname) VALUES('" + team["faction_id"].ToString() + "','" + team["leader"].ToString() + "','" + team["name"].ToString() + "')";
                                        command = new MySqlCommand(queryT, cnn2);
                                        reader2 = command.ExecuteReader();
                                        if (!reader2.HasRows)
                                        {
                                            if (team["faction_id"].ToString() != "bye")
                                            {
                                                command = new MySqlCommand(queryI, cnn);
                                                command.ExecuteNonQuery();
                                            }
                                        }
                                        reader2.Close();

                                        if (team["faction_id"].ToString() != "bye")
                                        {
                                            foreach (JToken player in team["roster"])
                                            {
                                                x++;
                                                string query2 = "Select * from player where player_id ='" + player["player_id"].ToString() + "'";
                                                string query = "";
                                                command = new MySqlCommand(query2, cnn2);
                                                reader1 = command.ExecuteReader();
                                                if (!reader1.HasRows)
                                                {
                                                    if (team["faction_id"].ToString() != "bye")
                                                    {
                                                        query = "INSERT INTO player (team_id,player_id,cup_id,playername,game_player_id) VALUES('" + team["faction_id"].ToString() + "','" + player["player_id"].ToString() + "','" + id + "','" + player["nickname"].ToString() + "','" + player["game_player_id"].ToString() + "')";
                                                    }
                                                }
                                                else
                                                {
                                                    query = "UPDATE player SET cup_id='" + id + "' WHERE player_id = '" + player["player_id"].ToString() + "'";
                                                }
                                                reader1.Close();
                                                command = new MySqlCommand(query, cnn);
                                                command.ExecuteNonQuery();
                                            }
                                        }
                                    }
                                }
                            }
                            reader4.Close();
                        }
                    }
                }
                cnn4.Close();
                cnn3.Close();
                cnn2.Close();
                cnn.Close();
            }
        }
    }
}
