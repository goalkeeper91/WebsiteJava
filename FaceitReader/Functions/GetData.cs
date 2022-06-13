using FaceitReader.Classes;
using FaceitReader.Database;
using MySql.Data.MySqlClient;
using Newtonsoft.Json;
using Newtonsoft.Json.Linq;
using System;
using System.Collections.Generic;
using System.IO;
using System.Net;

namespace FaceitReader.Functions
{
    internal class GetData
    {
        private static String Api = Globals.Api;

        public static List<FullPlaceList> getPlacementData(int limit, int match, string id, int cup)
        {
            var Purl = "https://open.faceit.com/data/v4/championships/" + id + "/results?offset=0&limit=" + limit;
            var Murl = "https://open.faceit.com/data/v4/championships/" + id + "/matches?type=past&offset=0&limit=" + match;
            var httpMRequest = (HttpWebRequest)WebRequest.Create(Murl);
            var httpPRequest = (HttpWebRequest)WebRequest.Create(Purl);

            List<Placement> placements = new List<Placement>();
            List<Team> teams = new List<Team>();
            List<Member> member = new List<Member>();
            List<Captain> captains = new List<Captain>();
            var list = new List<FullPlaceList>();

            int p = 0;
            int x = 1;
            int y = 0;

            httpPRequest.ContentType = "application/json; charset=utf-8";
            httpPRequest.Accept = "application/json";
            httpPRequest.Headers["Authorization"] = Api;
            var httpPResponse = (HttpWebResponse)httpPRequest.GetResponse();

            using (var streamReader = new StreamReader(httpPResponse.GetResponseStream()))
            {
                var result = streamReader.ReadToEnd();

                JObject json = JsonConvert.DeserializeObject<JObject>(result);

                foreach (var item in json["items"])
                {
                    foreach (var placement in item["placements"])
                    {
                        if (y == 0)
                        {
                            x = Int32.Parse(item["bounds"]["left"].ToString());
                            y++;
                        }
                        placements.Add(new Placement() { Place = x, TeamId = placement["id"].ToString(), TeamName = placement["name"].ToString() });
                        x++;
                    }
                }
            }

            httpMRequest.ContentType = "application/json; charset=utf-8";
            httpMRequest.Accept = "application/json";
            httpMRequest.Headers["Authorization"] = Api;
            var httpMResponse = (HttpWebResponse)httpMRequest.GetResponse();

            using (var streamReader = new StreamReader(httpMResponse.GetResponseStream()))
            {
                var result = streamReader.ReadToEnd();
                

                JObject json = JsonConvert.DeserializeObject<JObject>(result);

                foreach (JToken item in json["items"])
                {
                    foreach (JToken placement in item["teams"])
                    {
                        foreach (JToken team in placement.Children())
                        {
                            if (team["faction_id"].ToString() != "bye")
                            {
                                captains.Add(new Captain() { TeamId = team["faction_id"].ToString(), CaptainId = team["leader"].ToString() });
                                foreach (JToken player in team["roster"])
                                {
                                    member.Add(new Member() { TeamId = team["faction_id"].ToString(), PlayerId = player["player_id"].ToString(), PlayerName = player["nickname"].ToString() });
                                }
                            }
                        }
                    }
                }
            }

            foreach (var place in placements)
            {
                foreach (var captain in captains)
                {
                    if (place.TeamId == captain.TeamId)
                    {
                        foreach (var player in member)
                        {
                            if (player.PlayerId == captain.CaptainId)
                            {
                                teams.Add(new Team() { Place = place.Place, TeamName = place.TeamName, CaptainName = player.PlayerName });
                            }
                        }
                    }
                }
            }

            foreach (var team in teams)
            {
                if (team.Place != p)
                {
                    if (cup == 1)
                    {
                        if (team.Place == 1)
                        {
                            list.Add(new FullPlaceList() { Place = team.Place, TeamName = team.TeamName, CaptainName = team.CaptainName, Prize = "25 Punkte + 5x 25€" });
                            p = team.Place;
                        }
                        else if (team.Place == 2)
                        {
                            list.Add(new FullPlaceList() { Place = team.Place, TeamName = team.TeamName, CaptainName = team.CaptainName, Prize = "20 Punkte + 5x 15€" });
                            p = team.Place;
                        }
                        else if (team.Place == 3)
                        {
                            list.Add(new FullPlaceList() { Place = team.Place, TeamName = team.TeamName, CaptainName = team.CaptainName, Prize = "15 Punkte + 5x 10€" });
                            p = team.Place;
                        }
                        else if (team.Place == 4)
                        {
                            list.Add(new FullPlaceList() { Place = team.Place, TeamName = team.TeamName, CaptainName = team.CaptainName, Prize = "15 Punkte + 5x 10€" });
                            p = team.Place;
                        }
                        else if (team.Place == 5)
                        {
                            list.Add(new FullPlaceList() { Place = team.Place, TeamName = team.TeamName, CaptainName = team.CaptainName, Prize = "10 Punkte + 5x 5€" });
                            p = team.Place;
                        }
                        else if (team.Place == 6)
                        {
                            list.Add(new FullPlaceList() { Place = team.Place, TeamName = team.TeamName, CaptainName = team.CaptainName, Prize = "10 Punkte + 5x 5€" });
                            p = team.Place;
                        }
                        else if (team.Place == 7)
                        {
                            list.Add(new FullPlaceList() { Place = team.Place, TeamName = team.TeamName, CaptainName = team.CaptainName, Prize = "10 Punkte + 5x 5€" });
                            p = team.Place;
                        }
                        else if (team.Place == 8)
                        {
                            list.Add(new FullPlaceList() { Place = team.Place, TeamName = team.TeamName, CaptainName = team.CaptainName, Prize = "10 Punkte + 5x 5€" });
                            p = team.Place;
                        }
                        else if (team.Place > 8)
                        {
                            list.Add(new FullPlaceList() { Place = team.Place, TeamName = team.TeamName, CaptainName = team.CaptainName, Prize = "5x 2€" });
                            p = team.Place;
                        }
                    }
                    else if (cup == 2)
                    {
                        if (team.Place == 1)
                        {
                            list.Add(new FullPlaceList() { Place = team.Place, TeamName = team.TeamName, CaptainName = team.CaptainName, Prize = "5x 25€" });
                            p = team.Place;
                        }
                        else if (team.Place == 2)
                        {
                            list.Add(new FullPlaceList() { Place = team.Place, TeamName = team.TeamName, CaptainName = team.CaptainName, Prize = "5x 10€" });
                            p = team.Place;
                        }
                        else if (team.Place == 3)
                        {
                            list.Add(new FullPlaceList() { Place = team.Place, TeamName = team.TeamName, CaptainName = team.CaptainName, Prize = "5x 5€" });
                            p = team.Place;
                        }
                        else if (team.Place == 4)
                        {
                            list.Add(new FullPlaceList() { Place = team.Place, TeamName = team.TeamName, CaptainName = team.CaptainName, Prize = "5x 5€" });
                            p = team.Place;
                        }
                        else if (team.Place == 5)
                        {
                            list.Add(new FullPlaceList() { Place = team.Place, TeamName = team.TeamName, CaptainName = team.CaptainName, Prize = "5x 2€" });
                            p = team.Place;
                        }
                        else if (team.Place == 6)
                        {
                            list.Add(new FullPlaceList() { Place = team.Place, TeamName = team.TeamName, CaptainName = team.CaptainName, Prize = "5x 2€" });
                            p = team.Place;
                        }
                        else if (team.Place == 7)
                        {
                            list.Add(new FullPlaceList() { Place = team.Place, TeamName = team.TeamName, CaptainName = team.CaptainName, Prize = "5x 2€" });
                            p = team.Place;
                        }
                        else if (team.Place == 8)
                        {
                            list.Add(new FullPlaceList() { Place = team.Place, TeamName = team.TeamName, CaptainName = team.CaptainName, Prize = "5x 2€" });
                            p = team.Place;
                        }
                        else if (team.Place > 8)
                        {
                            list.Add(new FullPlaceList() { Place = team.Place, TeamName = team.TeamName, CaptainName = team.CaptainName, Prize = "5x 2€" });
                            p = team.Place;
                        }
                    }
                }
            }
            return list;
        }

        public static string getCupName(string id)
        {
            var name = "";
            var url = "https://open.faceit.com/data/v4/championships/" + id;
            var httpRequest = (HttpWebRequest)WebRequest.Create(url);

            httpRequest.ContentType = "application/json; charset=utf-8";
            httpRequest.Accept = "application/json";
            httpRequest.Headers["Authorization"] = Api;
            var httpResponse = (HttpWebResponse)httpRequest.GetResponse();
            using (var streamReader = new StreamReader(httpResponse.GetResponseStream()))
            {
                var result = streamReader.ReadToEnd();
                JObject json = JsonConvert.DeserializeObject<JObject>(result);
                name = json["name"].ToString();
            }

            return name;
        }
    }
}
