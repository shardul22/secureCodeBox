// SPDX-FileCopyrightText: the secureCodeBox authors
//
// SPDX-License-Identifier: Apache-2.0

import { set } from "lodash-es";

export async function handle({
  getFindings,
  updateFindings,
  attributeName = process.env["ATTRIBUTE_NAME"],
  attributeValue = process.env["ATTRIBUTE_VALUE"],
}) {
  const findings = await getFindings();

  const newFindings = findings.map((finding) => {
    set(finding, attributeName, attributeValue);
    return finding;
  });

  console.log(`Updated attributes on ${findings.length} findings`);

  await updateFindings(newFindings);
}
