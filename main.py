import argparse
import pandas as pd

parser = argparse.ArgumentParser()
parser.add_argument('--gitlog', type=str, help='Git log file to analyse')
args = parser.parse_args()

corp_domains = ["canonical.com", "collabora.com", "collabora.co.uk" "codethink.com", "endlessos.org", "intel.com", "openismus.com", "redhat.com", "ximian.com", "nvidia.com", "amd.com", "ubisoft.com", "microsoft.com", "apple.com", "huawei.com", "xilinx.com", "suse.de", "ibm.com", "linaro.org", "redhat.de", "codeweavers.com", "facebook.com", "netflix.com", "google.com", "xiaomi.com", "adobe.com", "docker.com", "oracle.com", "samsung.com", "suse.com", "qt.com", "qt.io", "qt-project.org", "ovi.com", "trolltech.com", "nokia.com", "mozilla.com", "nextcloud.com", "ubuntu.com", "collabora.co.uk", "suse.cz", "novell.com", "epicgames.com", "valvesoftware.com", "tensorflow.org", "swift-ci", "ibm.com", "fb.com", "twitter.com", "alibaba.com", "ge.com", "netscape.com", "ti.com", "citrix.com", "wolfsonmicro.com", "cisco.com", "fujitsu.com", "broadcom.com", "sgi.com", "hp.com", "atmel.com", "atheros.com", "nex.com", "coraid.com", "sun.com", "sony.co", "sony.com", "ntt.co", "ntt.com", "adaptec.com", "emulex.com", "analog.com", "vertias.com", "freescale.com", "qlogic.com", "toshiba.co", "toshiba.com", "arm.com", "marvell.com", "taobao.com", "micron.com", "hynix.com", "virtuozzo.com", "nxp.com", "linutronix.de", "free-electrons.com", "microsemi.com", "sang-engineering.com", "trendmicro.com", "rock-chips.com", "yandex-tem.ru", "altera.com", "alterra.com", "windriver.com", "synaptics.com", "codeaurora.org", "baylibre.com", "s-opensource.com", "savoirfairelinux.com", "mediatek.com", "lge.com", "lg.com", "renesas.com", "unisys.com", "qualcomm.com", "primarydata.com", "igalia.com", "aoyama.ac.jp", "unity.com", "shopify.com", "hulu.com", "rebertia.com", "kitware.com", "spotify.com", "wyeworks.com", "voormedia.com", "dio.jp", "zendesk.com", "slack-corp.com", "bqvision.com", "Obsidian.Systems"]

nextcloud_employees = ["nextcloud.com", "arthur-schiwon", "jus@bitgrid.net", "icewind.nl", "php.rio", "carlschwan.eu", "chrng8", "schilljs.com", "lchmn.me", "eneiluj", "vanpertsch", "artonge", "dependabot", "morrisjobke.de", "famdouma.nl", "winzerhof-wurst.at", "georgehrke.com", "rullzer", "danxuliu", "artificial-owl.com", "jancborchardt.net", "schiessle.org", "thomas.mueller", "owncloud.com", "owncloud-bot", "oparoz", "georgswebsite.de", "frank", "robin", "icewind1991", "karlitschek", "statuscode.ch"]

colnames = ["hash", "name", "email", "date", "subject"]
file = pd.read_csv(args.gitlog, parse_dates=[1], names=colnames)
file = file.dropna()

matches = file[file['email'].str.contains('|'.join(corp_domains))]
counter = len(matches)

percentage_corporate_contrib = (counter/len(file)) * 100

print("{}% of commits are from developers who are corporate associated".format(percentage_corporate_contrib))

no_date = file.drop('date', 1)
collapsed_by_dev_email = no_date.groupby(no_date['email'], as_index=False).sum()
matches = collapsed_by_dev_email[collapsed_by_dev_email['email'].str.contains('|'.join(corp_domains))]
counter = len(matches)
    
percentage_corporate = (counter/len(collapsed_by_dev_email)) * 100

print("{}% of developers are corporate associated".format(percentage_corporate))

to_date = file
to_date['date'] = pd.to_datetime(to_date['date'], utc=True)
to_date['date_y'] = to_date['date'].dt.to_period('Y')
by_year = to_date.groupby(to_date['date_y'], as_index=False).count()
by_year['yoy'] = by_year['email'].pct_change() * 100

print("Average YoY change in number of commits is {}%".format(by_year['yoy'].mean()))
