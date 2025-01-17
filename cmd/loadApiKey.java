import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.List;

public class LoadApiKey {

    public static String loadAPIKey() throws IOException {
        String homeDir = System.getProperty("user.home");
        Path filePath = Paths.get(homeDir, ".gocommit");
        List<String> lines;
        try {
            lines = Files.readAllLines(filePath);
        } catch (IOException e) {
            throw new IOException(
                "could not read .gocommit file. Did you run gocommit set-api --key hf_yourapikeyhere? ",
                e
            );
        }
        for (String line : lines) {
            if (line.startsWith("API_KEY=")) {
                return line.substring("API_KEY=".length());
            }
        }
        throw new IOException(
            "API_KEY not found in .gocommit file. Run gocommit set-api --key hf_yourapikeyhere."
        );
    }
}
