using System;
using System.Net.Http;
using System.Net.Http.Headers;
using System.Threading.Tasks;

namespace HcsDll
{
    public class Integration
    {
        static public async Task<string> PostDocument(string record, string documentType)
        {
            // Call asynchronous network methods in a try/catch block to handle exceptions.
            try
            {
                string HcsProxyUri = "";

                HttpClient client = new HttpClient();

                var stringContent = new StringContent(record);

                client.DefaultRequestHeaders.Accept.Add(new MediaTypeWithQualityHeaderValue("application/json"));

                var response = await client.PostAsync(HcsProxyUri + documentType, stringContent);

                return response.ToString();
            }
            catch (HttpRequestException e)
            {
                Console.WriteLine("\nException Caught!");
                Console.WriteLine("Message :{0} ", e.Message);
                return null;
            }
        }

    }
}
