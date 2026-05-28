import type { MemberInvite, SlackChannel, Team, TeamsChannel } from "@/types";

export const mockTeams: Team[] = [
  {
    id: "team-001",
    department: "Billing Team",
    categories: ["BILLING"],
    integrations: {
      slackWebhook: "https://hooks.slack.com/services/T00000000/B00000001/XXXXXXXX",
      slackChannelId: "C04AB1CD1EF",
      slackChannelName: "#billing-support",
      teamsWebhook: "https://outlook.office.com/webhook/billing",
      teamsChannelId: null,
      teamsChannelName: "Billing Support",
    },
    members: [
      {
        id: "mem-001",
        name: "Sarah Mitchell",
        emailAddress: "s.mitchell@acme.com",
        phoneNumber: "+44 20 7946 0101",
        role: "senior",
        integrations: { slackId: "U04AB1CD2EF", teamsId: "sarah.mitchell@acme.com" },
      },
      {
        id: "mem-002",
        name: "Connor Hayes",
        emailAddress: "c.hayes@acme.com",
        phoneNumber: "+44 20 7946 0102",
        role: "agent",
        integrations: { slackId: "U04AB1CD3GH", teamsId: "connor.hayes@acme.com" },
        inviteStatus: "pending",
        inviteId: "inv-002",
      },
    ],
    activeTicketCount: 7,
    createdAt: "2026-01-15T09:00:00Z",
  },
  {
    id: "team-002",
    department: "Technical Support",
    categories: ["TECHNICAL"],
    integrations: {
      slackWebhook: "https://hooks.slack.com/services/T00000000/B00000002/YYYYYYYY",
      slackChannelId: "C04AB1CD4IJ",
      slackChannelName: "#tech-support",
      teamsWebhook: null,
      teamsChannelId: null,
      teamsChannelName: null,
    },
    members: [
      {
        id: "mem-003",
        name: "James Okafor",
        emailAddress: "j.okafor@acme.com",
        phoneNumber: "+44 20 7946 0201",
        role: "senior",
        integrations: { slackId: "U04AB1CD4IJ", teamsId: "james.okafor@acme.com" },
      },
      {
        id: "mem-004",
        name: "Nina Rossi",
        emailAddress: "n.rossi@acme.com",
        phoneNumber: null,
        role: "agent",
        integrations: { slackId: "U04AB1CD5KL", teamsId: null },
      },
    ],
    activeTicketCount: 3,
    createdAt: "2026-01-15T09:00:00Z",
  },
  {
    id: "team-003",
    department: "Logistics & Delivery",
    categories: ["DELIVERY"],
    integrations: {
      slackWebhook: "https://hooks.slack.com/services/T00000000/B00000003/ZZZZZZZZ",
      slackChannelId: "C04AB1CD6MN",
      slackChannelName: "#logistics",
      teamsWebhook: "https://outlook.office.com/webhook/logistics",
      teamsChannelId: null,
      teamsChannelName: "Logistics Team",
    },
    members: [
      {
        id: "mem-005",
        name: "Priya Sharma",
        emailAddress: "p.sharma@acme.com",
        phoneNumber: "+44 20 7946 0301",
        role: "senior",
        integrations: { slackId: "U04AB1CD6MN", teamsId: "priya.sharma@acme.com" },
      },
      {
        id: "mem-006",
        name: "Luke Freeman",
        emailAddress: "l.freeman@acme.com",
        phoneNumber: "+44 20 7946 0302",
        role: "agent",
        integrations: { slackId: "U04AB1CD7OP", teamsId: null },
      },
    ],
    activeTicketCount: 2,
    createdAt: "2026-01-15T09:00:00Z",
  },
  {
    id: "team-004",
    department: "Customer Retention",
    categories: ["SUBSCRIPTIONS", "CANCELLATION"],
    integrations: {
      slackWebhook: "https://hooks.slack.com/services/T00000000/B00000004/AAAAAAAA",
      slackChannelId: "C04AB1CD8QR",
      slackChannelName: "#retention",
      teamsWebhook: null,
      teamsChannelId: null,
      teamsChannelName: null,
    },
    members: [
      {
        id: "mem-007",
        name: "Daniel Webb",
        emailAddress: "d.webb@acme.com",
        phoneNumber: "+44 20 7946 0401",
        role: "senior",
        integrations: { slackId: "U04AB1CD8QR", teamsId: "daniel.webb@acme.com" },
      },
      {
        id: "mem-008",
        name: "Fatima Al-Hassan",
        emailAddress: "f.alhassan@acme.com",
        phoneNumber: "+44 20 7946 0402",
        role: "agent",
        integrations: { slackId: "U04AB1CD9ST", teamsId: "fatima.alhassan@acme.com" },
      },
    ],
    activeTicketCount: 4,
    createdAt: "2026-01-15T09:00:00Z",
  },
  {
    id: "team-005",
    department: "Customer Success",
    categories: ["GENERAL", "ACCOUNT_ACCESS"],
    integrations: {
      slackWebhook: "https://hooks.slack.com/services/T00000000/B00000005/BBBBBBBB",
      slackChannelId: "C04AB1CDFGH",
      slackChannelName: "#customer-success",
      teamsWebhook: "https://outlook.office.com/webhook/customer-success",
      teamsChannelId: null,
      teamsChannelName: "Customer Success",
    },
    members: [
      {
        id: "mem-009",
        name: "Tom Baxter",
        emailAddress: "t.baxter@acme.com",
        phoneNumber: "+44 20 7946 0801",
        role: "senior",
        integrations: { slackId: "U04AB1CDF56", teamsId: "tom.baxter@acme.com" },
      },
      {
        id: "mem-010",
        name: "Isla Morgan",
        emailAddress: "i.morgan@acme.com",
        phoneNumber: "+44 20 7946 0802",
        role: "agent",
        integrations: { slackId: "U04AB1CDG78", teamsId: "isla.morgan@acme.com" },
        inviteStatus: "expired",
        inviteId: "inv-010",
      },
    ],
    activeTicketCount: 5,
    createdAt: "2026-01-15T09:00:00Z",
  },
];

export const mockInvites: MemberInvite[] = [
  {
    id: "inv-002",
    teamId: "team-001",
    email: "c.hayes@acme.com",
    name: "Connor Hayes",
    role: "agent",
    status: "pending",
    invitedAt: new Date(Date.now() - 9 * 60 * 60 * 1000).toISOString(),
    expiresAt: new Date(Date.now() + 39 * 60 * 60 * 1000).toISOString(),
    acceptedAt: null,
  },
];

export const mockSlackChannels: SlackChannel[] = [
  { id: "C01", name: "billing-support", isPrivate: false, memberCount: 12 },
  { id: "C02", name: "billing-general", isPrivate: false, memberCount: 34 },
  { id: "C03", name: "finance-team", isPrivate: true, memberCount: 8 },
];

export const mockTeamsChannels: TeamsChannel[] = [
  { id: "MS01", name: "Billing Support", teamName: "Finance", isPrivate: true },
  { id: "MS02", name: "General", teamName: "Finance", isPrivate: false },
  { id: "MS03", name: "Operations Alerts", teamName: "Operations", isPrivate: false },
];
