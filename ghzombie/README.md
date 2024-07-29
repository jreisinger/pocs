Ghzombie finds commits from deleted GitHub repositories. It's based on the idea about accessing deleted fork data from [this article](https://trufflesecurity.com/blog/anyone-can-access-deleted-and-private-repo-data-github).

```
$ wget https://data.gharchive.org/2024-07-25-{15..17}.json.gz # events from 3pm to 5pm UTC
$ gunzip *.gz
$ go run github.com/jreisinger/pocs/ghzombie@latest 2024-07-25-{15..17}.json 2> /dev/null
{
  "Deleted": "https://github.com/njms-lover/HANSAMAL-MD",
  "Commits": [
    {
      "Message": "Update config.cjs",
      "Url": "https://github.com/cobrs11/HANSAMAL-MD/commit/1eb95be27c0d3377418995aaadc0b617eec0d144"
    }
  ]
}
{
  "Deleted": "https://github.com/Mubashir6677/Suhail-Md",
  "Commits": [
    {
      "Message": "Create deploye.yml",
      "Url": "https://github.com/SuhailTechInfo/Suhail-Md/commit/f8c9815fb40062bb7cc19b9e73fceca900d509ba"
    },
    {
      "Message": "Update config.js",
      "Url": "https://github.com/SuhailTechInfo/Suhail-Md/commit/cebc53d9571e1a7b1ba758f982d4d65919dc1a1e"
    }
  ]
}
```