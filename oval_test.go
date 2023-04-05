package xmltree_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/pschou/go-xmltree"
)

var xccdfXML = `<?xml version='1.0' encoding='UTF-8'?>
<xccdf:Benchmark xmlns:xccdf="http://checklists.nist.gov/xccdf/1.1" xmlns:oval="http://oval.mitre.org/XMLSchema/oval-common-5" xmlns:oval-def="http://oval.mitre.org/XMLSchema/oval-definitions-5" xmlns:xhtml="http://www.w3.org/1999/xhtml" id="generated-xccdf" resolved="1">
  <xccdf:status>incomplete</xccdf:status>
  <xccdf:title>Red Hat Vulnerability Assessment for com.redhat.rhsa-all.xml</xccdf:title>
  <xccdf:description>This file has been automatically generated for purpose of vulnerability assessment of
            Red Hat products.</xccdf:description>
  <xccdf:rear-matter xml:lang="en-US">Red Hat and Red Hat Enterprise Linux are either registered trademarks or
            trademarks of Red Hat, Inc. in the United States and other countries. All other names are registered trademarks
            or trademarks of their respective companies.</xccdf:rear-matter>
  <xccdf:platform idref="cpe:/o:redhat:enterprise_linux"/>
  <xccdf:version time="2023-03-27T10:35:15">None, generated from OVAL file.</xccdf:version>
  <xccdf:Rule selected="true" id="oval-com.redhat.rhba-def-20070304" severity="high">
    <xccdf:title>RHBA-2007:0304: Updated kernel packages available for Red Hat Enterprise Linux 4 Update 5 (Important)</xccdf:title>
    <xccdf:description xml:lang="en-US">
      <xhtml:pre>Security Fix(es):

* kernel: tun: avoid double free in tun_free_netdev (CVE-2022-4744)

* ALSA: pcm: Move rwsem lock inside snd_ctl_elem_read to prevent UAF (CVE-2023-0266)

For more details about the security issue(s), including the impact, a CVSS score, acknowledgments, and other related information, refer to the CVE page(s) listed in the References section.</xhtml:pre>
    </xccdf:description>
    <xccdf:ident system="https://access.redhat.com/errata">RHSA-2023:1471</xccdf:ident>
    <xccdf:ident system="http://cve.mitre.org">CVE-2022-4744</xccdf:ident>
    <xccdf:ident system="http://cve.mitre.org">CVE-2023-0266</xccdf:ident>
    <xccdf:check system="http://oval.mitre.org/XMLSchema/oval-definitions-5">
      <xccdf:check-content-ref href="com.redhat.rhsa-all.xml" name="oval:com.redhat.rhsa:def:20231471"/>
    </xccdf:check>
  </xccdf:Rule>
</xccdf:Benchmark>`

var ovalXML = `<?xml version='1.0' encoding='UTF-8'?>
<oval_definitions xmlns="http://oval.mitre.org/XMLSchema/oval-definitions-5" xmlns:ind-def="http://oval.mitre.org/XMLSchema/oval-definitions-5#independent" xmlns:linux-def="http://oval.mitre.org/XMLSchema/oval-definitions-5#linux" xmlns:oval="http://oval.mitre.org/XMLSchema/oval-common-5" xmlns:oval-def="http://oval.mitre.org/XMLSchema/oval-definitions-5" xmlns:unix-def="http://oval.mitre.org/XMLSchema/oval-definitions-5#unix" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://oval.mitre.org/XMLSchema/oval-definitions-5#independent independent-definitions-schema.xsd http://oval.mitre.org/XMLSchema/oval-definitions-5#linux linux-definitions-schema.xsd http://oval.mitre.org/XMLSchema/oval-definitions-5#unix unix-definitions-schema.xsd http://oval.mitre.org/XMLSchema/oval-definitions-5 oval-definitions-schema.xsd http://oval.mitre.org/XMLSchema/oval-common-5 oval-common-schema.xsd">
  <generator>
    <oval:product_name>Debian</oval:product_name>
    <oval:schema_version>5.11.2</oval:schema_version>
    <oval:timestamp>2023-03-22T03:37:28.188-04:00</oval:timestamp>
  </generator>
  <definitions>
    <definition id="oval:org.debian:def:172501021062181471011117684875659626353" version="1" class="vulnerability">
      <metadata>
        <title>CVE-2022-3043 chromium</title>
        <affected family="unix">
          <platform>Debian GNU/Linux 11</platform>
          <product>chromium</product>
        </affected>
        <reference source="CVE" ref_id="CVE-2022-3043" ref_url="https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2022-3043"/>
        <description>security update</description>
        <debian>
          <dsa>DSA-5223</dsa>
          <moreinfo>Multiple security issues were discovered in Chromium, which could result in the execution of arbitrary code, denial of service or information disclosure.</moreinfo>
        </debian>
      </metadata>
      <criteria comment="Release section" operator="AND">
        <criterion test_ref="oval:org.debian.oval:tst:1" comment="Debian 11 is installed"/>
        <criteria comment="Architecture section" operator="OR">
          <criteria comment="Architecture independent section" operator="AND">
            <criterion test_ref="oval:org.debian.oval:tst:2" comment="all architecture"/>
            <criterion test_ref="oval:org.debian.oval:tst:24484" comment="chromium DPKG is earlier than 105.0.5195.52-1~deb11u1"/>
          </criteria>
        </criteria>
      </criteria>
    </definition>
  </definitions>
  <tests>
    <textfilecontent54_test id="oval:org.debian.oval:tst:1" version="1" check="all" check_existence="at_least_one_exists" comment="Debian GNU/Linux 11 is installed" xmlns="http://oval.mitre.org/XMLSchema/oval-definitions-5#independent">
      <object object_ref="oval:org.debian.oval:obj:1"/>
      <state state_ref="oval:org.debian.oval:ste:1"/>
    </textfilecontent54_test>
    <uname_test id="oval:org.debian.oval:tst:2" version="1" check="all" check_existence="at_least_one_exists" comment="Installed architecture is all" xmlns="http://oval.mitre.org/XMLSchema/oval-definitions-5#unix">
      <object object_ref="oval:org.debian.oval:obj:2"/>
    </uname_test>
    <dpkginfo_test id="oval:org.debian.oval:tst:24484" version="1" check="all" check_existence="at_least_one_exists" comment="chromium is earlier than 105.0.5195.52-1~deb11u1" xmlns="http://oval.mitre.org/XMLSchema/oval-definitions-5#linux">
      <object object_ref="oval:org.debian.oval:obj:1969"/>
      <state state_ref="oval:org.debian.oval:ste:16360"/>
    </dpkginfo_test>
  </tests>
  <objects>
    <textfilecontent54_object id="oval:org.debian.oval:obj:1" version="1" xmlns="http://oval.mitre.org/XMLSchema/oval-definitions-5#independent">
      <path>/etc</path>
      <filename>debian_version</filename>
      <pattern operation="pattern match">(\d+)\.\d</pattern>
      <instance datatype="int">1</instance>
    </textfilecontent54_object>
    <uname_object id="oval:org.debian.oval:obj:2" version="1" xmlns="http://oval.mitre.org/XMLSchema/oval-definitions-5#unix"/>
    <dpkginfo_object id="oval:org.debian.oval:obj:1969" version="1" xmlns="http://oval.mitre.org/XMLSchema/oval-definitions-5#linux">
      <name>chromium</name>
    </dpkginfo_object>
  </objects>
  <states>
    <textfilecontent54_state id="oval:org.debian.oval:ste:1" version="1" xmlns="http://oval.mitre.org/XMLSchema/oval-definitions-5#independent">
      <subexpression operation="equals">11</subexpression>
    </textfilecontent54_state>
    <dpkginfo_state id="oval:org.debian.oval:ste:16360" version="1" xmlns="http://oval.mitre.org/XMLSchema/oval-definitions-5#linux">
      <evr datatype="debian_evr_string" operation="less than">0:105.0.5195.52-1~deb11u1</evr>
    </dpkginfo_state>
  </states>
</oval_definitions>`

func TestRemoveLocal(t *testing.T) {
	root, err := xmltree.ParseXML(strings.NewReader(ovalXML))
	root.WalkFunc(func(el *xmltree.Element) error {
		el.RemoveLocalNS()
		return nil
	})
	//for i := range root.Children {
	//	trimSpace(&root.Children[i])
	//}
	//err = xmltree.EncodeIndent(os.Stdout, root, "", "\t")
	if err != nil {
		fmt.Println("err:", err)
	}
}

func TestSimplify(t *testing.T) {
	root, err := xmltree.ParseXML(strings.NewReader(xccdfXML))
	root.SimplifyNS()

	/*root.WalkAll(func(el *xmltree.Element) error {
		el.SimplifyNS()
		return nil
	})*/
	//for i := range root.Children {
	//	trimSpace(&root.Children[i])
	//}
	//err = xmltree.EncodeIndent(os.Stdout, root, "", "\t")
	if err != nil {
		fmt.Println("err:", err)
	}
	// Output:
	//
}
