using FaceitReader.Classes;
using Google.Apis.Auth.OAuth2;
using Google.Apis.Services;
using Google.Apis.Sheets.v4.Data;
using Google.Apis.Sheets.v4;
using Google.Apis.Util.Store;
using Newtonsoft.Json;
using Newtonsoft.Json.Linq;
using System;
using System.Collections.Generic;
using System.Globalization;
using System.IO;
using System.Linq;
using System.Net;
using System.Threading;
using System.Windows.Forms;
using NPOI.SS.UserModel;
using NPOI.SS.Formula.Functions;

namespace FaceitReader.Functions
{
    internal class GetData
    {
        private static String Api = Globals.Api;
        static readonly string[] Scopes = { SheetsService.Scope.Spreadsheets };
        static readonly string ApplicationName = "Faceit Reader";

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
                                teams.Add(new Team() { Place = place.Place, TeamName = place.TeamName, TeamId = place.TeamId, CaptainName = player.PlayerName });
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
                            list.Add(new FullPlaceList() { Place = team.Place, TeamName = team.TeamName, TeamId = team.TeamId, CaptainName = team.CaptainName, Prize = "25 Punkte + 5x 25€" });
                            p = team.Place;
                        }
                        else if (team.Place == 2)
                        {
                            list.Add(new FullPlaceList() { Place = team.Place, TeamName = team.TeamName, TeamId = team.TeamId, CaptainName = team.CaptainName, Prize = "20 Punkte + 5x 15€" });
                            p = team.Place;
                        }
                        else if (team.Place == 3)
                        {
                            list.Add(new FullPlaceList() { Place = team.Place, TeamName = team.TeamName, TeamId = team.TeamId, CaptainName = team.CaptainName, Prize = "15 Punkte + 5x 10€" });
                            p = team.Place;
                        }
                        else if (team.Place == 4)
                        {
                            list.Add(new FullPlaceList() { Place = team.Place, TeamName = team.TeamName, TeamId = team.TeamId, CaptainName = team.CaptainName, Prize = "15 Punkte + 5x 10€" });
                            p = team.Place;
                        }
                        else if (team.Place == 5)
                        {
                            list.Add(new FullPlaceList() { Place = team.Place, TeamName = team.TeamName, TeamId = team.TeamId, CaptainName = team.CaptainName, Prize = "10 Punkte + 5x 5€" });
                            p = team.Place;
                        }
                        else if (team.Place == 6)
                        {
                            list.Add(new FullPlaceList() { Place = team.Place, TeamName = team.TeamName, TeamId = team.TeamId, CaptainName = team.CaptainName, Prize = "10 Punkte + 5x 5€" });
                            p = team.Place;
                        }
                        else if (team.Place == 7)
                        {
                            list.Add(new FullPlaceList() { Place = team.Place, TeamName = team.TeamName, TeamId = team.TeamId, CaptainName = team.CaptainName, Prize = "10 Punkte + 5x 5€" });
                            p = team.Place;
                        }
                        else if (team.Place == 8)
                        {
                            list.Add(new FullPlaceList() { Place = team.Place, TeamName = team.TeamName, TeamId = team.TeamId, CaptainName = team.CaptainName, Prize = "10 Punkte + 5x 5€" });
                            p = team.Place;
                        }
                        else if (team.Place > 8)
                        {
                            list.Add(new FullPlaceList() { Place = team.Place, TeamName = team.TeamName, TeamId = team.TeamId, CaptainName = team.CaptainName, Prize = "5x 2€" });
                            p = team.Place;
                        }
                    }
                    else if (cup == 2)
                    {
                        if (team.Place == 1)
                        {
                            list.Add(new FullPlaceList() { Place = team.Place, TeamName = team.TeamName, TeamId = team.TeamId, CaptainName = team.CaptainName, Prize = "5x 25€" });
                            p = team.Place;
                        }
                        else if (team.Place == 2)
                        {
                            list.Add(new FullPlaceList() { Place = team.Place, TeamName = team.TeamName, TeamId = team.TeamId, CaptainName = team.CaptainName, Prize = "5x 10€" });
                            p = team.Place;
                        }
                        else if (team.Place == 3)
                        {
                            list.Add(new FullPlaceList() { Place = team.Place, TeamName = team.TeamName, TeamId = team.TeamId, CaptainName = team.CaptainName, Prize = "5x 5€" });
                            p = team.Place;
                        }
                        else if (team.Place == 4)
                        {
                            list.Add(new FullPlaceList() { Place = team.Place, TeamName = team.TeamName, TeamId = team.TeamId, CaptainName = team.CaptainName, Prize = "5x 5€" });
                            p = team.Place;
                        }
                        else if (team.Place == 5)
                        {
                            list.Add(new FullPlaceList() { Place = team.Place, TeamName = team.TeamName, TeamId = team.TeamId, CaptainName = team.CaptainName, Prize = "5x 2€" });
                            p = team.Place;
                        }
                        else if (team.Place == 6)
                        {
                            list.Add(new FullPlaceList() { Place = team.Place, TeamName = team.TeamName, TeamId = team.TeamId, CaptainName = team.CaptainName, Prize = "5x 2€" });
                            p = team.Place;
                        }
                        else if (team.Place == 7)
                        {
                            list.Add(new FullPlaceList() { Place = team.Place, TeamName = team.TeamName, TeamId = team.TeamId, CaptainName = team.CaptainName, Prize = "5x 2€" });
                            p = team.Place;
                        }
                        else if (team.Place == 8)
                        {
                            list.Add(new FullPlaceList() { Place = team.Place, TeamName = team.TeamName, TeamId = team.TeamId, CaptainName = team.CaptainName, Prize = "5x 2€" });
                            p = team.Place;
                        }
                        else if (team.Place > 8)
                        {
                            list.Add(new FullPlaceList() { Place = team.Place, TeamName = team.TeamName, TeamId = team.TeamId, CaptainName = team.CaptainName, Prize = "5x 2€" });
                            p = team.Place;
                        }
                    }
                }
            }
            return list;
        }
        public static string getDate(string id)
        {
            DateTime date = new DateTime(1970, 1, 1, 0, 0, 0, 0, DateTimeKind.Utc);
            string datestring;
            var Purl = "https://open.faceit.com/data/v4/championships/" + id;
            var httpPRequest = (HttpWebRequest)WebRequest.Create(Purl);

            httpPRequest.ContentType = "application/json; charset=utf-8";
            httpPRequest.Accept = "application/json";
            httpPRequest.Headers["Authorization"] = Api;
            var httpPResponse = (HttpWebResponse)httpPRequest.GetResponse();

            using (var streamReader = new StreamReader(httpPResponse.GetResponseStream()))
            {
                var result = streamReader.ReadToEnd();

                JObject json = JsonConvert.DeserializeObject<JObject>(result);
                long unix = (long)json["championship_start"];
                date = date.AddMilliseconds(unix).ToLocalTime();
            }
            datestring = date.ToString("dd.MM.yyyy", CultureInfo.InvariantCulture);
            return datestring;
        }
        public static List<Tuple<string, string, string>> getTeamDetails(string teamname)
        {
            List<Tuple<string, string, string>> list = new List<Tuple<string, string, string>>();
            var url = "https://open.faceit.com/data/v4/teams/" + teamname;
            var httpRequest = (HttpWebRequest)WebRequest.Create(url);

            httpRequest.ContentType = "application/json; charset=utf-8";
            httpRequest.Accept = "application/json";
            httpRequest.Headers["Authorization"] = Api;
            var httpResponse = (HttpWebResponse)httpRequest.GetResponse();
            using (var streamReader = new StreamReader(httpResponse.GetResponseStream()))
            {
                var result = streamReader.ReadToEnd();
                JObject json = JsonConvert.DeserializeObject<JObject>(result);

                if (json != null)
                {
                    list.Add(new Tuple<string, string, string>(json["team_id"].ToString(), json["name"].ToString(), "https://www.faceit.com/de/teams/" + teamname));
                }
                if (list.Count() == 0)
                {
                    list.Add(new Tuple<string, string, string>("No Team found", "", ""));
                }
            }
            return list;
        }
        public static void CreateEntry(DataGridView dt1, string id, string rang)
        {
            var valueRange = new ValueRange();
            UserCredential credential;

            using (var stream =
                new FileStream("Resources/credentials.json", FileMode.Open, FileAccess.Read))
            {
                string credPath = "token.json";
                credential = GoogleWebAuthorizationBroker.AuthorizeAsync(
                    GoogleClientSecrets.FromStream(stream).Secrets,
                    Scopes,
                    "user",
                    CancellationToken.None,
                    new FileDataStore(credPath, true)).Result;
            }
            var service = new SheetsService(new BaseClientService.Initializer()
            {
                HttpClientInitializer = credential,
                ApplicationName = ApplicationName,
            });

            var objectList = new List<object>();
            try
            {
                foreach (DataGridViewRow row in dt1.Rows)
                {
                    objectList.Add(row.Cells[2].Value.ToString() + " via Faceit an " + row.Cells[4].Value.ToString());
                    valueRange.Values = new List<IList<object>> { objectList };
                }

            }
            catch
            {

            }
            valueRange.MajorDimension = "COLUMNS";
            var appendRequest = service.Spreadsheets.Values.Update(valueRange, id, rang);
            appendRequest.ValueInputOption = SpreadsheetsResource.ValuesResource.UpdateRequest.ValueInputOptionEnum.USERENTERED;
            appendRequest.Execute();
        }
        public static void updateSpreadsheetCells(string spreadsheetId, string sheetName, int columnId, int rowStart, int rowEnd, string? team, string? captain)
        {
            UserCredential credential;

            using (var stream =
                new FileStream("Resources/credentials.json", FileMode.Open, FileAccess.Read))
            {
                string credPath = "token.json";
                credential = GoogleWebAuthorizationBroker.AuthorizeAsync(
                    GoogleClientSecrets.FromStream(stream).Secrets,
                    Scopes,
                    "user",
                    CancellationToken.None,
                    new FileDataStore(credPath, true)).Result;
            }
            var service = new SheetsService(new BaseClientService.Initializer()
            {
                HttpClientInitializer = credential,
                ApplicationName = ApplicationName,
            });

            //get sheet id by sheet name
            Spreadsheet spr = service.Spreadsheets.Get(spreadsheetId).Execute();
            Sheet sh = spr.Sheets.Where(s => s.Properties.Title == sheetName).FirstOrDefault();
            int sheetId = (int)sh.Properties.SheetId;

            //define cell color
            var userEnteredFormat = new CellFormat()
            {
                BackgroundColor = new Color()
                {
                    Blue = 0,
                    Red = 0,
                    Green = 1,
                    Alpha = (float)1.0
                },
                Borders = new Borders()
                {
                    Bottom = new Border() { Style = "solid"},
                    Top = new Border() { Style = "solid" },
                    Right = new Border() { Style = "solid" },
                    Left = new Border() { Style = "solid" }
                }
            };

            var userEnteredValue = new ExtendedValue()
            {
                StringValue = team + " via Faceit an " + captain + " am " + DateTime.Now.ToString("dd.MM.yyyy")
            };

            BatchUpdateSpreadsheetRequest bussr = new BatchUpdateSpreadsheetRequest();
            var updateCellsRequest = new Request();
            //create the update request for cells from the first row
            if (team != null)
            {
                updateCellsRequest = new Request()
                {
                    RepeatCell = new RepeatCellRequest()
                    {
                        Range = new GridRange()
                        {
                            SheetId = sheetId,
                            StartColumnIndex = columnId,
                            StartRowIndex = rowStart,
                            EndColumnIndex = columnId + 1,
                            EndRowIndex = rowEnd + 1
                        },
                        Cell = new CellData()
                        {
                            UserEnteredFormat = userEnteredFormat,
                            UserEnteredValue = userEnteredValue
                        },
                        Fields = "*"
                    }
                };
            }
            else
            {
                updateCellsRequest = new Request()
                {
                    RepeatCell = new RepeatCellRequest()
                    {
                        Range = new GridRange()
                        {
                            SheetId = sheetId,
                            StartColumnIndex = columnId,
                            StartRowIndex = rowStart,
                            EndColumnIndex = columnId + 1,
                            EndRowIndex = rowEnd + 1
                        },
                        Cell = new CellData()
                        {
                            UserEnteredFormat = userEnteredFormat,
                        },
                        Fields = "UserEnteredFormat(BackgroundColor)"
                    }
                };
            }
            
            bussr.Requests = new List<Request>();
            bussr.Requests.Add(updateCellsRequest);
            var bur = service.Spreadsheets.BatchUpdate(bussr, spreadsheetId);
            bur.Execute();
        }
    }
}
