<div align="center">
<pre style="line-height: 1;">
<span style="color: #ff5faf;">█▀▄▀█ ████▄ ████▄ ██▄     ▄▀  ▄█    ▄▄▄▄▀</span>
<span style="color: #d75fff;">█ █ █ █   █ █   █ █  █  ▄▀    ██ ▀▀▀ █   </span>
<span style="color: #af5fff;">█ ▄ █ █   █ █   █ █   █ █ ▀▄  ██     █   </span>
<span style="color: #875fff;">█   █ ▀████ ▀████ █  █  █   █ ▐█    █    </span>
<span style="color: #5f5fff;">   █              ███▀   ███   ▐    ▀    </span>
<span style="color: #0000ff;">  ▀                                      </span>
</pre>
</div>

# moodgit

a simple CLI tool to log and track your mood through the command line. moodgit helps you maintain a personal mood journal, allowing you to track your emotional state with intensity levels, descriptive messages, and tags to build insights into your mood patterns over time.

## features

- 🎭 **multiple mood types**: track various emotions including happy, sad, angry, anxious, excited, calm, stressed, tired, and neutral
- 📊 **intensity scale**: rate your mood intensity from 0-10 for more detailed tracking
- 📝 **custom messages**: add descriptive messages to provide context for your mood entries
- 🏷️ **tagging system**: organize your entries with custom tags for better categorization
- 📚 **mood history**: view your mood logs in chronological order
- ✏️ **entry amendment**: modify your last mood entry if needed

## installation

### prerequisites

- Go 1.25.1 or later

### build from source

1. clone the repository:

   ```bash
   git clone https://github.com/niitish/moodgit.git
   cd moodgit
   ```

2. build the project:

   ```bash
   go build -o moodgit
   ```

3. (optional) move the binary to your PATH:

   ```bash
   # on Linux/macOS
   sudo mv moodgit /usr/local/bin/

   # on Windows, add the directory to your PATH environment variable
   ```

## quick start

1. **initialize your mood repository**:

   ```bash
   moodgit init
   ```

2. **add your first mood entry**:

   ```bash
   moodgit add -i 10 -o happy -m "yo! this works!" -t achievement
   ```

3. **view your mood history**:
   ```bash
   moodgit log
   ```

## data storage

moodgit stores your mood data locally in a SQLite database located at `~/.moodgit/moodgit.db`. your data remains private and is never transmitted anywhere.

## contributing

contributions are welcome! please feel free to submit a pull request. for major changes, please open an issue first to discuss what you would like to change.
