﻿// Copyright (c) 2021 EPAM Systems
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

using System;

namespace Epam.FixAntenna.NetCore.Message
{
	/// <summary>
	/// Thrown when an invalid message detected. For
	/// example, sequence is to low or senderCompId is differ than expected.
	/// </summary>
	internal class InvalidMessageException : ArgumentException
	{
		private readonly bool _critical;
		private readonly FixMessage _invalidMessage;

		public InvalidMessageException(FixMessage invalidMessage) : this(invalidMessage, "Invalid message")
		{
		}

		public InvalidMessageException(FixMessage invalidMessage, string cause) : this(invalidMessage, cause, false)
		{
		}

		public InvalidMessageException(FixMessage invalidMessage, string cause, bool critical) : base(cause)
		{
			_invalidMessage = invalidMessage;
			_critical = critical;
		}

		public override string Message => base.Message + ": " + _invalidMessage.ToPrintableString();

		public virtual FixMessage GetInvalidMessage()
		{
			return _invalidMessage;
		}

		public override string ToString()
		{
			return Message;
		}

		public virtual bool IsCritical()
		{
			return _critical;
		}
	}
}